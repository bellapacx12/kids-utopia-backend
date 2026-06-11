package model

import "time"

type BookVariant struct {
	ID        string    `db:"id"`
	BookID    string    `db:"book_id"`

	Language  string    `db:"language"` // en | am | or | ti
	Title     string    `db:"title"`

	FileURL   string    `db:"file_url"`
	FileKey   string    `db:"file_key"`

	Status    string    `db:"status"`   // uploading | processing | ready | failed
	Progress  int       `db:"progress"` // 0–100
	Error     string    `db:"error"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}