package milestones

import (
	"context"
	"log"

	"github.com/bellapacx/kids-utopia/internal/gamification/model"
	"github.com/bellapacx/kids-utopia/internal/gamification/rules"
)

type Repository interface {
	HasAwarded(ctx context.Context, childID, milestoneID, bookID string) (bool, error)
	Award(ctx context.Context, award *model.MilestoneAward) error
	InsertEvent(ctx context.Context, e *model.MilestoneEvent) (string, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Process(
	ctx context.Context,
	event rules.Event,
	rewards []rules.Reward,
) error {

	// 1. create event FIRST
	eventID, err := s.repo.InsertEvent(ctx, &model.MilestoneEvent{
		ChildID: event.ChildID,
		BookID:  event.BookID,
		Type:    "milestone_trigger",
		Title:   "Processing milestones",
	Meta: map[string]any{
	"book_id":      event.BookID,
},
	})

	
	if err != nil {
		return err
	}

	for _, r := range rewards {

		if r.Type != rules.RewardTypeBadge &&
			r.Type != rules.RewardTypeTheme {
			continue
		}

		milestoneID, ok := r.Meta["milestone_id"].(string)
		if !ok || milestoneID == "" {
			continue
		}
        
		exists, err := s.repo.HasAwarded(
			ctx,
			event.ChildID,
			milestoneID,
			event.BookID,
		)
		if err != nil {
			return err
		}

		if exists {
			continue
		}
        log.Println("AWARD INSERT DEBUG")
log.Println("childID:", event.ChildID)
log.Println("bookID:", event.BookID)
log.Println("eventID:", eventID)
log.Println("milestoneID:", milestoneID)
		// 2. award linked to event
		err = s.repo.Award(ctx, &model.MilestoneAward{
			ChildID:     event.ChildID,
			BookID:      event.BookID,
			MilestoneID: milestoneID,
			RewardType:  r.Type,
			RewardRef:   nil,
			EventID:     eventID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func generateTitle(id string) string {
	switch id {
	case "first_page":
		return "🎉 First page read!"
	case "book_completed":
		return "📚 Book completed!"
	case "streak_7":
		return "🔥 7-day streak!"
	default:
		return "🏆 Achievement unlocked!"
	}
}