package model

import "time"

type SessionStart struct {
	ChildID   string
	SessionID string
	StartedAt time.Time
}
type ProgressEvent struct {
	ChildID   string
	BookID    string
	Page      int
	RecordedAt time.Time
}
type SessionEnd struct {
	ChildID   string
	SessionID string
	Duration  int
	EndedAt   time.Time
}