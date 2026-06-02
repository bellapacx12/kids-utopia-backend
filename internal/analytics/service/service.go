package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/bellapacx/kids-utopia/internal/analytics/model"
	"github.com/bellapacx/kids-utopia/internal/analytics/repository"
	"github.com/bellapacx/kids-utopia/internal/events"
	sessionrepo "github.com/bellapacx/kids-utopia/internal/reader_session/repository"
	streakrepo "github.com/bellapacx/kids-utopia/internal/streak/repository"
)

type Service struct {
	repo *repository.Repo
	streakRepo  streakrepo.StreakRepository
	sessionRepo  sessionrepo.SessionRepository
}

func New(r *repository.Repo,  streakRepo streakrepo.StreakRepository, sessionRepo sessionrepo.SessionRepository) *Service {
	return &Service{repo: r, streakRepo: streakRepo, sessionRepo: sessionRepo}
}
func (s *Service) ProcessMessage(ctx context.Context, msg string) error {

	var event events.Event

	if err := json.Unmarshal([]byte(msg), &event); err != nil {
		return err
	}

	dbEvent := model.Event{
		EventID:   event.EventID,
		Type:      string(event.Type),
		UserID:    event.UserID,
		ChildID:   event.ChildID,
		BookID:    event.BookID,
		SessionID: event.SessionID,
		Page:      event.Page,
		Timestamp: event.Timestamp,
	}
	log.Printf("📊 ANALYTICS EVENT INSERT: %+v", dbEvent)
	if err := s.repo.Insert(ctx, dbEvent); err != nil {
	log.Printf("❌ ANALYTICS INSERT FAILED: %v", err)
	return err
}

log.Printf("✅ ANALYTICS INSERT OK: %s", dbEvent.EventID)

	return s.repo.Insert(ctx, dbEvent)
}
func (s *Service) GetAnalytics(ctx context.Context, childID string) (map[string]any, error) {

	// 1. aggregate analytics events
	stats, err := s.repo.GetChildStats(ctx, childID)
	if err != nil {
		return nil, err
	}

	// 2. get streak
	streak, err := s.streakRepo.Get(ctx, childID)
	if err != nil {
		// if no streak yet → fallback
		streak = nil
	}
	timeSpent, err := s.sessionRepo.GetTotalReadingTime(ctx, childID)
if err != nil {
	return nil, err
}

	result := map[string]any{
		"child_id": childID,

		"total_sessions": stats.TotalSessions,
		"total_pages": stats.TotalPages,
		"total_reading_time_seconds": timeSpent,

	}

	if streak != nil {
		result["current_streak"] = streak.CurrentStreak
		result["longest_streak"] = streak.LongestStreak
	} else {
		result["current_streak"] = 0
		result["longest_streak"] = 0
	}

	return result, nil
}
