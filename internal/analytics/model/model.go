package model

import "time"

type Event struct {
	EventID   string    `json:"event_id"`
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	ChildID   string    `json:"child_id"`
	BookID    string    `json:"book_id"`
	SessionID string    `json:"session_id"`
	Page      int       `json:"page"`
	Timestamp time.Time `json:"timestamp"`
}