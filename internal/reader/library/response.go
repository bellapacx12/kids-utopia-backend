package library

import "time"

type ContinueReadingItem struct {
	BookID          string    `json:"book_id"`
	Title           string    `json:"title"`
	CoverURL        string    `json:"cover_url"`
	CurrentPage     int       `json:"current_page"`
	ProgressPercent int       `json:"progress_percent"`
	LastReadAt      time.Time `json:"last_read_at"`
}

type ContinueReadingResponse struct {
	Items []ContinueReadingItem `json:"items"`
}