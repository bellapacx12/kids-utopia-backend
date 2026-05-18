package auth

import (
	"context"
	"time"

	redisClient "github.com/bellapacx/kids-utopia/pkg/redis"
)
func StoreRefreshSession(userID string, token string) error {
	key := "session:" + userID + ":" + token

	return redisClient.Client.Set(
		context.Background(),
		key,
		"active",
		7*24*time.Hour,
	).Err()
}