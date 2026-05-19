package auth

import (
	"context"
	"time"

	redis "github.com/bellapacx/kids-utopia/pkg/redis"
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
func DeleteRefreshSession(token string) error {

	return redis.Client.Del(
		context.Background(),
		"refresh:"+token,
	).Err()
}