package xp

import (
	"github.com/bellapacx/kids-utopia/internal/gamification/rules"
)

type PageProgressRule struct{}

func (r PageProgressRule) Match(event rules.Event, state rules.State) bool {
	return event.Type == "progress.updated"
}

func (r PageProgressRule) Execute(
	event rules.Event,
	state rules.State,
) ([]rules.Reward, error) {

	// ignore invalid events
	if event.Page <= 0 {
		return nil, nil
	}

	// no backward or duplicate progress
	if event.Page <= event.PreviousPage {
		return nil, nil
	}

	xp := event.Page - event.PreviousPage

	return []rules.Reward{
		{
			Type:  rules.RewardTypeXP,
			Value: xp,
			Meta: map[string]any{
				"from": event.PreviousPage,
				"to":   event.Page,
			},
		},
	}, nil
}