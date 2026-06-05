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

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	analyticsrepo "github.com/bellapacx/kids-utopia/internal/analytics/repository"
	analyticssvc "github.com/bellapacx/kids-utopia/internal/analytics/service"
	"github.com/bellapacx/kids-utopia/internal/books/events"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	gamificationrules "github.com/bellapacx/kids-utopia/internal/gamification/rules"
	"github.com/bellapacx/kids-utopia/internal/worker"

	streakrepo "github.com/bellapacx/kids-utopia/internal/streak/repository"
	streaksvc "github.com/bellapacx/kids-utopia/internal/streak/service"

	// READER SESSION
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
	"github.com/bellapacx/kids-utopia/pkg/sqs"
	"github.com/bellapacx/kids-utopia/pkg/storage"
)

func main() {
	cfg := config.Load()

	// =========================
	// DATABASE
	// =========================
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

	// =========================
	// REPOSITORY
	// =========================
	bookPagesRepo := repository.NewBookPagesRepository(database.DB)

	// =========================
	// STORAGE
	// =========================
	st, err := storage.NewS3Storage(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
		cfg.S3PublicURL,
	)
	if err != nil {
		log.Fatal("storage init errorrrr:", err)
	}

	// =========================
	// SQS
	// =========================
	queue, err := sqs.New(cfg.SQSQueueURL, cfg.AWSRegion)
	if err != nil {
		log.Fatal("sqs init error:", err)
	}

	log.Println("🚀 Worker running (SQS consumer)")
	sessionRepo := sessionrepo.New(database.DB)
streakRepo := streakrepo.New(database.DB)
streakservice := streaksvc.New(streakRepo)
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
	streakservice,
	progressService,
	themesService,
)
	// =========================
	// CONTEXT / SHUTDOWN
	// =========================
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Println("🛑 shutdown signal received")
		cancel()
	}()

	// =========================
	// WORKER POOL
	// =========================
	const workerCount = 3
	jobs := make(chan types.Message, 20)

	for i := 0; i < workerCount; i++ {
		go workerLoop(ctx,i, jobs, queue, st, bookPagesRepo, analyticsService, gamificationService)
	}

	// =========================
	// POLLER LOOP
	// =========================
	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 worker shutting down")
			close(jobs)
			return

		default:
			messages, err := queue.Receive()
			if err != nil {
				log.Println("receive error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			for _, msg := range messages {
				jobs <- msg
			}
		}
	}
}

func workerLoop(
	ctx context.Context,
	id int,
	jobs <-chan types.Message,
	queue *sqs.Client,
	st storage.Storage,
	repo repository.BookPagesRepository,
	analyticsService *analyticssvc.Service,
	gamificationService *gamificationsvc.Service,
) {
	for msg := range jobs {

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("🔥 panic recovered:", r)
				}
			}()

			// =========================
			// STEP 1: decode base event (lightweight routing)
			// =========================
			var base struct {
				Type string `json:"type"`
			}

			if err := json.Unmarshal([]byte(*msg.Body), &base); err != nil {
				log.Println("❌ invalid base event:", err)
				return
			}

			log.Printf("📩 worker %d received event type: %s", id, base.Type)

			// =========================
			// STEP 2: route by type
			// =========================
			switch base.Type {

			// =========================
			// BOOK UPLOADED (existing logic untouched)
			// =========================
			case "book.uploaded":

				var event events.BookUploadedEvent

				if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
					log.Println("❌ book event decode failed:", err)
					return
				}

				log.Printf("📘 worker %d processing book: %s", id, event.BookID)

				if err := worker.ProcessBook(event, st, repo); err != nil {
					log.Printf("❌ worker %d failed book %s: %v", id, event.BookID, err)
					return
				}

				log.Printf("✅ worker %d completed book: %s", id, event.BookID)

			// =========================
			// PROGRESS EVENT (placeholder)
			// =========================
			case "progress.updated":
                log.Println("update event")
				if err := analyticsService.ProcessMessage(ctx, *msg.Body); err != nil {
		log.Printf("❌ analytics insert failed (progress.updated): %v", err)
		return
	}
	          
	          var event appEvents.Event

	if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
		log.Printf("❌ progress decode failed: %v", err)
		return
	}

err := gamificationService.ProcessEvent(
	ctx,
	gamificationrules.Event{
		Type:      string(event.Type),
		ChildID:   event.ChildID,
		SessionID: event.SessionID,
		BookID:    event.BookID,
		Page:      event.Page,
		PreviousPage: event.PreviousPage,
		EventID:   event.EventID,
		TotalPages : event.TotalPages,
	},
)

if err != nil {
	log.Printf("❌ gamification failed: %v", err)
	return
}

	log.Printf(
		"🎮 XP awarded child=%s xp=1 event=%s",
		event.ChildID,
		event.EventID,
	)

			// =========================
			// SESSION EVENTS (placeholder)
			// =========================
			case "session.started":
				
	               if err := analyticsService.ProcessMessage(ctx, *msg.Body); err != nil {
		log.Printf("❌ analytics insert failed: %v", err)
	}

			case "session.ended":
			if err := analyticsService.ProcessMessage(ctx, *msg.Body); err != nil {
		log.Printf("❌ analytics insert failed: %v", err)
	}        
			// =========================
			// UNKNOWN
			// =========================
			default:
				log.Printf("⚠️ unknown event type: %s", base.Type)
			}

			// =========================
			// DELETE MESSAGE ONLY AFTER SUCCESS ROUTE
			// =========================
			if err := queue.Delete(*msg.ReceiptHandle); err != nil {
				log.Println("❌ delete error:", err)
				return
			}

		}()
	}
}