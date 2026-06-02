package events

import "time"

// =========================
// BASE EVENT (used everywhere: engine, SQS, worker)
// =========================

type Event struct {
	EventID   string    `json:"event_id"`
	Type      EventType `json:"type"`

	UserID    string    `json:"user_id"`
	ChildID   string    `json:"child_id"`
	BookID    string    `json:"book_id"`
	SessionID string    `json:"session_id,omitempty"`

	Page      int       `json:"page"`
	Timestamp time.Time `json:"timestamp"`
}