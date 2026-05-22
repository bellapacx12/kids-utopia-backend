package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"

	"github.com/bellapacx/kids-utopia/internal/subscriptions/dto"
	"github.com/bellapacx/kids-utopia/internal/subscriptions/model"
	"github.com/bellapacx/kids-utopia/internal/subscriptions/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}
func (s *Service) Create(ctx context.Context, userID string, req dto.CreateSubscriptionRequest) error {

	var duration time.Duration

	switch req.Plan {
	case "monthly":
		duration = 30 * 24 * time.Hour
	case "yearly":
		duration = 365 * 24 * time.Hour
	default:
		duration = 30 * 24 * time.Hour
	}

	now := time.Now()
	end := now.Add(duration)

	sub := model.Subscription{
		ID:        uuid.NewString(),
		UserID:    userID,
		Plan:      req.Plan,
		Status:    "active",
		StartDate: now,
		EndDate:   &end,
	}

	return s.repo.Create(ctx, sub)
}
func (s *Service) HasActive(
	ctx context.Context,
	userID string,
) (bool, error) {

	sub, err := s.repo.GetActiveByUser(ctx, userID)

	if err != nil {
		// treat any "no rows" variant safely
		if errors.Is(err, pgx.ErrNoRows) ||
			strings.Contains(err.Error(), "no rows in result set") {
			return false, nil
		}
		return false, err
	}

	return sub != nil, nil
}