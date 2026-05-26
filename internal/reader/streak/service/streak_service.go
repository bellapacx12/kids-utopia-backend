package service

import (
	"context"
	"time"

	"github.com/bellapacx/kids-utopia/internal/reader/streak/model"
	"github.com/bellapacx/kids-utopia/internal/reader/streak/repository"
)

type StreakService struct {
	repo repository.StreakRepository
}

func New(repo repository.StreakRepository) *StreakService {
	return &StreakService{repo: repo}
}
func (s *StreakService) UpdateStreak(ctx context.Context, childID string) error {

	today := time.Now().UTC().Truncate(24 * time.Hour)

streak, err := s.repo.Get(ctx, childID)

if err != nil {

	newStreak := &model.ReadingStreak{
		ChildID:       childID,
		CurrentStreak: 1,
		LongestStreak: 1,
		LastReadDate:  today,
		UpdatedAt:     time.Now(),
	}

	return s.repo.Create(ctx, newStreak)
}
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
streak.UpdatedAt = time.Now()

return s.repo.Update(ctx, streak)
}
