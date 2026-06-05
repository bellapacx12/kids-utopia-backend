package model

import "time"

type ChildXP struct {
	ChildID   string    `json:"child_id"`
	TotalXP   int       `json:"total_xp"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type XPTransaction struct {
	ID        string    `json:"id"`
	ChildID   string    `json:"child_id"`
	Source    string    `json:"source"`
	SourceID  string    `json:"source_id"`
	XPAmount  int       `json:"xp_amount"`
	CreatedAt time.Time `json:"created_at"`
}
// MilestoneAward is the system-level record (idempotency / safety layer)
type MilestoneAward struct {
	ID          string
	ChildID     string
	BookID      string
	MilestoneID string
	RewardType  string
	RewardRef   *string
	EventID     string
	CreatedAt   time.Time
}
type MilestoneEvent struct {
	ID          string
	ChildID     string
	BookID      string
	Type        string

	Title       string
	Description string

	Meta        map[string]any

	CreatedAt   time.Time
}