package model

import "time"

type BookProgress struct {
	ID string `db:"id" json:"id"`

	ChildID string `db:"child_id" json:"child_id"`
	BookID  string `db:"book_id" json:"book_id"`

	CurrentPage int  `db:"current_page" json:"current_page"`
	Completed   bool `db:"completed" json:"completed"`

	ProgressPercent int `db:"progress_percent" json:"progress_percent"`

	LastReadAt time.Time `db:"last_read_at" json:"last_read_at"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}