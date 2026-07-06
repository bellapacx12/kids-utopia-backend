package app

import (
	"fmt"
	"log"

	"github.com/bellapacx/kids-utopia/internal/reader/events"
	"github.com/bellapacx/kids-utopia/pkg/config"
	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/kafka"
	"github.com/bellapacx/kids-utopia/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	redisClient "github.com/bellapacx/kids-utopia/pkg/redis"
	sqsClient "github.com/bellapacx/kids-utopia/pkg/sqs"

	"github.com/bellapacx/kids-utopia/pkg/storage"
)

type Container struct {
	Config  *config.Config
	Storage *storage.S3Storage
	DB      *pgxpool.Pool
	Queue   *sqsClient.Client
	ReaderEventsBus *events.Bus
	KafkaProducer   *kafka.Producer
}

func NewContainer() *Container {

	cfg := config.Load()

	logger.Init()

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

	redisClient.Connect(cfg.RedisURL)

	storageClient, err := storage.NewS3Storage(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
		cfg.S3PublicURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	queue, err := sqsClient.New(
		cfg.SQSQueueURL,
		cfg.AWSRegion,
	)
	if err != nil {
		log.Fatal(err)
	}
	bus := events.NewBus(queue)
     
    kafkaClient, err := kafka.New(kafka.Config{
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

producer := kafka.NewProducer(kafkaClient)

	return &Container{
		Config:  cfg,
		DB:      database.DB,
		Storage: storageClient,
		Queue:   queue,
		KafkaProducer:   producer,
		ReaderEventsBus: bus,
	}
}