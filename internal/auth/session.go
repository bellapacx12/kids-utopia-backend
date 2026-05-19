package auth

import (
	"context"
	"time"

	redis "github.com/bellapacx/kids-utopia/pkg/redis"
)
func StoreRefreshSession(userID string, token string) error {

	key := "refresh:" + token

	return redis.Client.Set(
		context.Background(),
		key,
		userID,
		7*24*time.Hour,
	).Err()
}
func DeleteRefreshSession(token string) error {

	return redis.Client.Del(
		context.Background(),
		"refresh:"+token,
	).Err()
}