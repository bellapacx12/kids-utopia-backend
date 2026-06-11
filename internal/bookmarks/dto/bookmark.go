package dto

type BookmarkDTO struct {
	ID        string `json:"id"`
	BookID    string `json:"book_id"`
	Page      int    `json:"page"`

	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
}
