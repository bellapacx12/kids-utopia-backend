package rules

// Event is the normalized gamification input
	type Event struct {
		Type      string
		ChildID   string
		SessionID string
		BookID    string
		Page      int
		PreviousPage int `json:"previous_page"`
		TotalPages int     `json:"total_pages"`
		EventID   string
	}

	// State represents current known user gamification state
	type State struct {
		LastPage int
		TotalXP  int
		Streak   int
		Level    int
		TotalPages int // REQUIRED for completion logic
		BookSeen      bool
	BookCompleted bool
	}

	// Reward is output of rules engine
	type Reward struct {
		Type  string                 // xp, badge, streak, theme
		Value int
		Meta  map[string]any
	}

	// Rule is the core interface all gamification rules implement
	type Rule interface {
		Match(event Event, state State) bool
		Execute(event Event, state State) ([]Reward, error)
	}
	const (
		RewardTypeXP     = "xp"
		RewardTypeBadge  = "badge"
		RewardTypeTheme  = "theme"
		RewardTypeStreak = "streak"
	)