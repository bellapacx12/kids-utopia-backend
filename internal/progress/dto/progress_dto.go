package dto

type UpdateProgressRequest struct {
	ChildID   string `json:"child_id"`
	BookID    string `json:"book_id"`
	Page      int    `json:"page"`
	TotalPages int   `json:"total_pages"`
}
type ProgressResponse struct {
	ChildID string `json:"child_id"`
	BookID  string `json:"book_id"`

	CurrentPage int `json:"current_page"`
	ProgressPercent int `json:"progress_percent"`

	Completed bool `json:"completed"`
}