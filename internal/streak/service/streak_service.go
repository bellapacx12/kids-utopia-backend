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
func (s *StreakService) ProcessSessionEnd(
	ctx context.Context,
	childID string,
	timestamp time.Time,
) error {

	streak, err := s.repo.Get(ctx, childID)

	if err != nil {
		// if not found → create
		streak = &model.ReadingStreak{
			ChildID:       childID,
			CurrentStreak: 1,
			LongestStreak: 1,
			LastReadDate:  timestamp,
		}
		return s.repo.Create(ctx, streak)
	}

	// streak logic (simple MVP)
	if sameDay(streak.LastReadDate, timestamp) {
		return nil
	}

	if yesterday(streak.LastReadDate, timestamp) {
		streak.CurrentStreak++
	} else {
		streak.CurrentStreak = 1
	}

	if streak.CurrentStreak > streak.LongestStreak {
		streak.LongestStreak = streak.CurrentStreak
	}

	streak.LastReadDate = timestamp

	return s.repo.Update(ctx, streak)
}
func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.YearDay() == b.YearDay()
}

func yesterday(last, now time.Time) bool {
	return now.Year() == last.Year() &&
		now.YearDay()-last.YearDay() == 1
}