package otp

import (
	"context"
	"time"

	redis "github.com/bellapacx/kids-utopia/pkg/redis"
)
type Service struct {
	router *Router
	
}

func NewService(router *Router) *Service {
	return &Service{
		router: router,
		
	}
}


func (s *Service) Verify(to string, code string) bool {

	val, err := redis.Client.Get(
		context.Background(),
		"otp:"+to,
	).Result()

	if err != nil {
		return false
	}

	if val != code {
		return false
	}

	redis.Client.Del(context.Background(), "otp:"+to)

	return true
}
func (s *Service) Send(to string, code string) error {

	err := redis.Client.Set(
		context.Background(),
		"otp:"+to,
		code,
		5*time.Minute,
	).Err()

	if err != nil {
		return err
	}

	subject := "Your Kids Utopia OTP Code "

	content := "<h2>Your OTP Code</h2><p><b>" + code + "</b></p><p>Expires in 5 minutes</p>"

	return s.router.Send(to, subject, content)
}