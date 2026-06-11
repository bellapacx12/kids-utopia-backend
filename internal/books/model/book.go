package model

import "time"

type Book struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Author      string    `db:"author"`
	CoverURL    string    `db:"cover_url"`
    Titles      []string
	Status      string    `db:"status"` // draft, processing, ready, failed

	AccessType  string    `db:"access_type"` // free, premium

	// 📊 Recommendation + UX fields
	AgeMin      int       `db:"age_min"`
	AgeMax      int       `db:"age_max"`
	Language    string    `db:"language"`
	Category    string    `db:"category"`
	PopularityScore int   `db:"popularity_score"`
    Progress int `db:"progress"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
