package model

import "time"

type ReadingSession struct {
	ID string `db:"id" json:"id"`

	UserID  string `db:"user_id" json:"user_id"`
	ChildID string `db:"child_id" json:"child_id"`
	BookID  string `db:"book_id" json:"book_id"`

	StartedAt time.Time  `db:"started_at" json:"started_at"`
	EndedAt   *time.Time `db:"ended_at" json:"ended_at"`

	DurationSeconds int `db:"duration_seconds" json:"duration_seconds"`

	StartPage int `db:"start_page" json:"start_page"`
	EndPage   *int `db:"end_page" json:"end_page"`

	Completed bool `db:"completed" json:"completed"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}