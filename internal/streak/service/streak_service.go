package service

import (
	"context"
	"time"

	"github.com/bellapacx/kids-utopia/internal/streak/model"
	"github.com/bellapacx/kids-utopia/internal/streak/repository"
)

type StreakService struct {
	repo repository.StreakRepository
}

func New(repo repository.StreakRepository) *StreakService {
	return &StreakService{repo: repo}
}
func (s *StreakService) UpdateStreak(ctx context.Context, childID string) error {

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	streak, err := s.repo.Get(ctx, childID)

	// create if missing
	if err != nil {
		newStreak := &model.ReadingStreak{
			ChildID:       childID,
			CurrentStreak: 1,
			LongestStreak: 1,
			LastReadDate:  today,
			UpdatedAt:     now,
		}
		return s.repo.Create(ctx, newStreak)
	}

	// idempotency guard
	if streak.LastReadDate.Equal(today) {
		return nil
	}

	yesterday := today.AddDate(0, 0, -1)

	if streak.LastReadDate.Equal(yesterday) {
		streak.CurrentStreak++
	} else {
		streak.CurrentStreak = 1
	}

	if streak.CurrentStreak > streak.LongestStreak {
		streak.LongestStreak = streak.CurrentStreak
	}

	streak.LastReadDate = today
	streak.UpdatedAt = now

	return s.repo.Update(ctx, streak)
}
func (s *StreakService) GetStreak(ctx context.Context, childID string) (*model.ReadingStreak, error) {
	return s.repo.Get(ctx, childID)
}
