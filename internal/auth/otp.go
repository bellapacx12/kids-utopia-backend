package auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	redis "github.com/bellapacx/kids-utopia/pkg/redis"
)

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(999999))
}

func StoreOTP(identifier string, code string) {
	redis.Client.Set(
		context.Background(),
		"otp:"+identifier,
		code,
		5*time.Minute,
	)
}

func VerifyOTP(identifier string, code string) bool {
	val, err := redis.Client.Get(context.Background(), "otp:"+identifier).Result()
	if err != nil {
		return false
	}
	return val == code
}