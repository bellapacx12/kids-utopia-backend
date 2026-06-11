package model

import "time"

type BookmarkDetail struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	BookTitle string    `json:"book_title"`
    CoverImage string   `json:"cover_image"`
	Page      int       `json:"page"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`

	CreatedAt time.Time `json:"created_at"`
}