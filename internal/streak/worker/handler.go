package worker

import (
	"context"
	"encoding/json"

	analyticsmodel "github.com/bellapacx/kids-utopia/internal/analytics/model"
	streaksvc "github.com/bellapacx/kids-utopia/internal/streak/service"
)

type Service struct {
	streak *streaksvc.StreakService
}

func New(s *streaksvc.StreakService) *Service {
	return &Service{streak: s}
}

func (w *Service) ProcessMessage(ctx context.Context, msg string) error {

	var event analyticsmodel.Event
	if err := json.Unmarshal([]byte(msg), &event); err != nil {
		return err
	}

	// ONLY SESSION END affects streak
	if event.Type != "session_ended" {
		return nil
	}

	return w.streak.UpdateStreak(ctx, event.ChildID)
}