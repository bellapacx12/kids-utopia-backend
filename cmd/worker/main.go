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

	"github.com/bellapacx/kids-utopia/internal/books/events"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/internal/worker"

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
		go workerLoop(i, jobs, queue, st, bookPagesRepo)
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
	id int,
	jobs <-chan types.Message,
	queue *sqs.Client,
	st storage.Storage,
	repo repository.BookPagesRepository,
) {
	for msg := range jobs {

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("🔥 panic recovered:", r)
				}
			}()

			var event events.BookUploadedEvent

			if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
				log.Println("invalid event:", err)
				return
			}

			log.Printf("📘 worker %d processing book: %s", id, event.BookID)

			if err := worker.ProcessBook(event, st, repo); err != nil {
				log.Printf("❌ worker %d failed book %s: %v", id, event.BookID, err)
				return
			}

			if err := queue.Delete(*msg.ReceiptHandle); err != nil {
				log.Println("delete error:", err)
				return
			}

			log.Printf("✅ worker %d completed book: %s", id, event.BookID)
		}()
	}
}