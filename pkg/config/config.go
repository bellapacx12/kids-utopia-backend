package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisURL string

	JWTSecret string

	KafkaBroker string

	// ✅ ADD THESE (S3 / MinIO)
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	 S3PublicURL string

	 AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	SESFromEmail       string

	SQSQueueURL string
	SendGridAPIKey string
	FromEmail string
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
	KafkaUsername string
	KafkaPassword string
	KafkaCAFile  string
	
}

func Load() *Config {
		_ = godotenv.Load()
	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		RedisURL: os.Getenv("REDIS_URL"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		KafkaBroker: os.Getenv("KAFKA_BROKER"),

		S3Endpoint:  os.Getenv("S3_ENDPOINT"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3PublicURL: os.Getenv("S3_PUBLIC_URL"),

		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),

		SESFromEmail: os.Getenv("SES_FROM_EMAIL"),
		SQSQueueURL:  os.Getenv("SQS_QUEUE_URL"),
		SendGridAPIKey : os.Getenv("SendGridAPIKey"),
		FromEmail : os.Getenv("FROM_EMAIL"),
        
		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
KafkaUsername: os.Getenv("KAFKA_USERNAME"),
KafkaPassword: os.Getenv("KAFKA_PASSWORD"),
KafkaCAFile:  os.Getenv("KAFKA_CA_FILE"),
	}
}