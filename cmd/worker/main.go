package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	analyticsmodel "github.com/bellapacx/kids-utopia/internal/analytics/model"
	analyticsrepo "github.com/bellapacx/kids-utopia/internal/analytics/repository"
	analyticssvc "github.com/bellapacx/kids-utopia/internal/analytics/service"

	bookevents "github.com/bellapacx/kids-utopia/internal/books/events"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/internal/worker"

	gamificationrules "github.com/bellapacx/kids-utopia/internal/gamification/rules"

	streakrepo "github.com/bellapacx/kids-utopia/internal/streak/repository"
	streaksvc "github.com/bellapacx/kids-utopia/internal/streak/service"

	sessionrepo "github.com/bellapacx/kids-utopia/internal/reader_session/repository"

	gamificationrepo "github.com/bellapacx/kids-utopia/internal/gamification/repository"
	gamificationsvc "github.com/bellapacx/kids-utopia/internal/gamification/service"

	milestones "github.com/bellapacx/kids-utopia/internal/gamification/milestones"
	milestonerepo "github.com/bellapacx/kids-utopia/internal/gamification/milestones/repository"

	appEvents "github.com/bellapacx/kids-utopia/internal/events"
	themes "github.com/bellapacx/kids-utopia/internal/gamification/themes"

	progressrepo "github.com/bellapacx/kids-utopia/internal/progress/repository"
	progressservice "github.com/bellapacx/kids-utopia/internal/progress/service"

	"github.com/bellapacx/kids-utopia/pkg/config"
	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/kafka"
	"github.com/bellapacx/kids-utopia/pkg/storage"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	cfg := config.Load()

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	database.Connect(dbURL)
	log.Println("✅ PostgreSQL connected")

	bookPagesRepo := repository.NewBookPagesRepository(database.DB)

	st, err := storage.NewS3Storage(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
		cfg.S3PublicURL,
	)
	if err != nil {
		log.Fatal("storage init error:", err)
	}

	// =========================
	// KAFKA
	// =========================
	client, err := kafka.New(kafka.Config{
		Brokers:  cfg.KafkaBrokers,
		Username: cfg.KafkaUsername,
		Password: cfg.KafkaPassword,
		CAFile:   cfg.KafkaCAFile,
		Topic:    cfg.KafkaTopic,
		GroupID:  cfg.KafkaGroupID,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("🧠 Kafka connected. Brokers=%v Topic=%s Group=%s",
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaGroupID,
	)

	consumer := kafka.NewConsumer(client)

	// =========================
	// SERVICES
	// =========================
	sessionRepo := sessionrepo.New(database.DB)
	streakRepo := streakrepo.New(database.DB)
	streakService := streaksvc.New(streakRepo)

	analyticsRepo := analyticsrepo.New(database.DB)
	analyticsService := analyticssvc.New(analyticsRepo, streakRepo, sessionRepo)

	gamificationRepo := gamificationrepo.New(database.DB)
	milestoneRepo := milestonerepo.New(database.DB)
	milestoneService := milestones.New(milestoneRepo)

	progressRepo := progressrepo.NewProgressRepository(database.DB)
	progressService := progressservice.NewProgressService(progressRepo)

	themesRepo := themes.NewRepository(database.DB)
	themesService := themes.New(themesRepo)

	gamificationService := gamificationsvc.New(
		gamificationRepo,
		milestoneService,
		streakService,
		progressService,
		themesService,
	)

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Println("🛑 shutdown signal received")
		cancel()
	}()

	const workerCount = 3
	jobs := make(chan *kgo.Record, 50)

	for i := 0; i < workerCount; i++ {
		go workerLoop(
			ctx,
			i,
			jobs,
			st,
			bookPagesRepo,
			analyticsService,
			gamificationService,
			streakService,
			consumer,
		)
	}

	log.Println("🚀 Worker running (Kafka consumer)")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 worker shutting down")
			close(jobs)
			return

		default:
			records, err := consumer.Poll(ctx)
			if err != nil {
				log.Println("❌ kafka poll error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			for _, r := range records {

				log.Printf("📨 RECEIVED MESSAGE: topic=%s partition=%d offset=%d key=%s",
					r.Topic,
					r.Partition,
					r.Offset,
					string(r.Key),
				)

				log.Printf("📦 RAW VALUE: %s", string(r.Value))

				jobs <- r
			}
		}
	}
}

func workerLoop(
	ctx context.Context,
	id int,
	jobs <-chan *kgo.Record,
	st storage.Storage,
	repo repository.BookPagesRepository,
	analyticsService *analyticssvc.Service,
	gamificationService *gamificationsvc.Service,
	streakService *streaksvc.StreakService,
	consumer *kafka.Consumer,
) {
	for msg := range jobs {

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("🔥 panic recovered:", r)
				}
			}()

			log.Printf("⚙️ Worker %d processing offset=%d topic=%s",
				id, msg.Offset, msg.Topic,
			)

			body := msg.Value

			var base struct {
				Type string `json:"type"`
			}

			if err := json.Unmarshal(body, &base); err != nil {
				log.Println("❌ invalid base event:", err)
				return
			}

			log.Printf("📩 worker %d event type=%s payload=%s", id, base.Type, string(body))

			switch base.Type {

			case "book.uploaded":

				var event bookevents.BookVariantUploaded
				if err := json.Unmarshal(body, &event); err != nil {
					log.Println("❌ decode failed:", err)
					return
				}

				if err := worker.ProcessBook(event, st, repo); err != nil {
					log.Println("❌ worker failed:", err)
					return
				}

			case "progress.updated":
				_ = analyticsService.ProcessMessage(ctx, string(body))

				var event appEvents.Event
				_ = json.Unmarshal(body, &event)

				_ = gamificationService.ProcessEvent(ctx, gamificationrules.Event{
					Type:         string(event.Type),
					ChildID:      event.ChildID,
					SessionID:    event.SessionID,
					BookID:       event.BookID,
					Page:         event.Page,
					PreviousPage: event.PreviousPage,
					EventID:      event.EventID,
					TotalPages:   event.TotalPages,
				})

			case "session.started":
				_ = analyticsService.ProcessMessage(ctx, string(body))

			case "session.ended":
				_ = analyticsService.ProcessMessage(ctx, string(body))

				var event analyticsmodel.Event
				_ = json.Unmarshal(body, &event)

				_ = streakService.UpdateStreak(ctx, event.ChildID)

			default:
				log.Printf("⚠️ unknown event type=%s raw=%s", base.Type, string(body))
			}

			if err := consumer.Commit(ctx, msg); err != nil {
				log.Println("❌ kafka commit error:", err)
			}
		}()
	}
}