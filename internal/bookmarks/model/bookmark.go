package model

import "time"

type Bookmark struct {
	ID        string    `json:"id"`
	ChildID   string    `json:"child_id"`
	BookID    string    `json:"book_id"`
	Page      int       `json:"page"`
	CreatedAt time.Time `json:"created_at"`
}

