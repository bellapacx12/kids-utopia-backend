package service

import (
	"context"
	"strings"

	"github.com/bellapacx/kids-utopia/internal/books/model"
	"github.com/bellapacx/kids-utopia/internal/subscriptions/service"
)
type AccessResult struct {
	Allowed bool
	Preview bool
}
type Service struct {
	subscriptionService *service.Service
}

func New(sub *service.Service) *Service {
	return &Service{
		subscriptionService: sub,
	}
}
func (s *Service) CanAccessBook(
	ctx context.Context,
	userID string,
	book *model.Book,
) (allowed bool, preview bool, err error) {

	// FREE BOOK → full access
	if strings.ToLower(book.AccessType) == "free" {
		return true, false, nil
	}

	// check subscription
	hasSub, err := s.subscriptionService.HasActive(ctx, userID)
	if err != nil {
		return false, false, err
	}

	// PREMIUM + NO SUB → PREVIEW
	if !hasSub {
		return false, true, nil
	}

	// PREMIUM + SUB → FULL ACCESS
	return true, false, nil
}