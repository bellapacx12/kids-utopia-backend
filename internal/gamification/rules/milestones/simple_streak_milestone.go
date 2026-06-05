package milestones

import "github.com/bellapacx/kids-utopia/internal/gamification/rules"

type StreakMilestone struct{}

func (m StreakMilestone) Match(event rules.Event, state rules.State) bool {
	return state.Streak > 0 && state.Streak%7 == 0
}

func (m StreakMilestone) Execute(
	event rules.Event,
	state rules.State,
) ([]rules.Reward, error) {

	// IMPORTANT: avoid duplicate badge spam
	if state.Streak == 0 {
		return nil, nil
	}

	return []rules.Reward{
		{
			Type:  rules.RewardTypeBadge,
			Value: 1,
			Meta: map[string]any{
				"milestone_id": "streak_7",
				"streak_value": state.Streak,
			},
		},
	}, nil
}