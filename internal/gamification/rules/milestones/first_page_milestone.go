package milestones

import "github.com/bellapacx/kids-utopia/internal/gamification/rules"

type FirstPageMilestone struct{}

func (r FirstPageMilestone) Match(event rules.Event, state rules.State) bool {
	return event.Type == "progress.updated" &&
		event.Page == 1 &&
		!state.BookSeen
}
func (r FirstPageMilestone) Execute(event rules.Event, state rules.State) ([]rules.Reward, error) {

	return []rules.Reward{
		{
			Type: rules.RewardTypeBadge,
			Value: 0,
			Meta: map[string]any{
				"milestone_id": "first_page",
				"book_id":      event.BookID,
			},
		},
		{
			Type:  rules.RewardTypeXP,
			Value: 5,
			Meta: map[string]any{
				"reason": "first_page_bonus",
			},
		},
	}, nil
}