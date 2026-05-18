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
	log.Println("✅ Worker PostgreSQL connected")

	bookPagesRepo := repository.NewBookPagesRepository(database.DB)

	st, err := storage.NewS3Storage(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
		cfg.S3PublicURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	queue, err := sqs.New(cfg.SQSQueueURL, cfg.AWSRegion)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 Worker started... listening to SQS")

	// ================================
	// CONTEXT (FIXED)
	// ================================
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Println("🛑 shutdown signal received")
		cancel()
	}()
		const workerCount = 3
	jobs := make(chan types.Message, 20)

	for i := 0; i < workerCount; i++ {
		go startWorker(i, jobs, queue, st, bookPagesRepo)
	}
		for {
		select {

		case <-ctx.Done():
			log.Println("🛑 worker shutting down...")
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
func startWorker(
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
					log.Println("🔥 worker panic recovered:", r)
				}
			}()

			var event events.BookUploadedEvent

			if err := json.Unmarshal([]byte(*msg.Body), &event); err != nil {
				log.Println("invalid event:", err)
				return
			}

			log.Printf("📘 worker %d processing book: %s", id, event.BookID)

			err := worker.ProcessBook(event, st, repo)
			if err != nil {
				log.Printf("❌ worker %d failed book %s: %v", id, event.BookID, err)
				return
			}

			err = queue.Delete(*msg.ReceiptHandle)
			if err != nil {
				log.Println("delete error:", err)
				return
			}

			log.Printf("✅ worker %d completed book: %s", id, event.BookID)
		}()
	}
}