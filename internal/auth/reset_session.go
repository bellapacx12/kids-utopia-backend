package auth

import (
	"context"
	"time"

	"github.com/bellapacx/kids-utopia/pkg/redis"
)

func StoreResetSession(identifier string) error {

	return redis.Client.Set(
		context.Background(),
		"reset_session:"+identifier,
		"active",
		10*time.Minute,
	).Err()
}
func ValidateResetSession(identifier string) (bool, error) {

	val, err := redis.Client.Get(
		context.Background(),
		"reset_session:"+identifier,
	).Result()

	if err != nil {
		return false, nil // not found = invalid
	}

	return val == "active", nil
}
func DeleteResetSession(identifier string) error {

	return redis.Client.Del(
		context.Background(),
		"reset_session:"+identifier,
	).Err()
}