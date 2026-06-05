package milestones

import "github.com/bellapacx/kids-utopia/internal/gamification/rules"

type BookCompletedMilestone struct{}
func (r BookCompletedMilestone) Match(event rules.Event, state rules.State) bool {
	return event.Type == "progress.updated" &&
		event.Page > 0 &&
		event.Page == event.TotalPages
}
func (r BookCompletedMilestone) Execute(event rules.Event, state rules.State) ([]rules.Reward, error) {

	return []rules.Reward{
		{
			Type: rules.RewardTypeBadge,
			Meta: map[string]any{
				"milestone_id": "book_completed",
				"book_id":      event.BookID,
			},
		},
		{
			Type: rules.RewardTypeTheme,
			Meta: map[string]any{
				"theme_id": "space",
			},
		},
		{
			Type:  rules.RewardTypeXP,
			Value: 20,
		},
	}, nil
}