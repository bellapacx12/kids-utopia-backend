package service

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/books/model"
	"github.com/bellapacx/kids-utopia/internal/subscriptions/service"
)

type Service struct {
	subscriptionService *service.Service
}

func New(sub *service.Service) *Service {
	return &Service{
		subscriptionService: sub,
	}
}
func (s *Service) CanAccessBook(ctx context.Context, userID string, book *model.Book) (bool, error) {

	// FREE BOOK → always allowed
	if book.AccessType == "free" {
		return true, nil
	}

	// PREMIUM BOOK → must be logged in
	if userID == "" {
		return false, nil
	}

	// check subscription
	hasSub, err := s.subscriptionService.HasActive(ctx, userID)
	if err != nil {
		return false, err
	}

	return hasSub, nil
}