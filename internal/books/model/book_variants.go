package model

import "time"

type BookVariant struct {
	ID        string    `db:"id"`
	BookID    string    `db:"book_id"`
	Language  string    `db:"language"` // en, am, sw, etc
	Title     string    `db:"title"`
	Summary   string    `db:"summary"`
	CreatedAt time.Time `db:"created_at"`
}