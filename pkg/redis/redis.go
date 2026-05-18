package redis

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func Connect(redisURL string) {

	opt, err := goredis.ParseURL(redisURL)
	if err != nil {
		log.Fatal(err)
	}

	Client = goredis.NewClient(opt)

	_, err = Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Redis connected")
}