package registry

import (
	"github.com/bellapacx/kids-utopia/internal/gamification/rules"
	"github.com/bellapacx/kids-utopia/internal/gamification/rules/milestones"
	"github.com/bellapacx/kids-utopia/internal/gamification/rules/xp"
)

func NewRules() []rules.Rule {
	return []rules.Rule{
		xp.PageProgressRule{},
		milestones.FirstPageMilestone{},
	milestones.BookCompletedMilestone{},
	milestones.StreakMilestone{},
	}
}