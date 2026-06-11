package dto

type Recommendation struct {
	BookID   string `json:"book_id"`
	Title    string `json:"title"`
	CoverURL string `json:"cover_url"`
	Score    int    `json:"score"`
	Reason   string `json:"reason"`
}