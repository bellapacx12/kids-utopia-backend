package model

import "time"

type Book struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description  string    `db:"description"`
	Author      string    `db:"author"`
	CoverURL    string    `db:"cover_url"`
	Status      string    `db:"status"` // draft, processing, ready, failed
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	AccessType  string    `db:"access_type"`
}