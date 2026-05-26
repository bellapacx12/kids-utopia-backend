package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/reader/streak/model"
)

type StreakRepository interface {
	Get(ctx context.Context, childID string) (*model.ReadingStreak, error)
	Create(ctx context.Context, s *model.ReadingStreak) error
	Update(ctx context.Context, s *model.ReadingStreak) error
}