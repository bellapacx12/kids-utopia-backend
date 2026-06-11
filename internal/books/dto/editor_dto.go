package dto

type EditorPageDTO struct {
	PageNumber int    `json:"page_number"`
	Content    string `json:"content"`
	ImageKey   string `json:"image_key"`
	ImageURL   string `json:"image_url"`
}

type EditorResponse struct {
	BookID string          `json:"book_id"`
	Status string          `json:"status"`
	Progress int           `json:"progress"`
	Pages  []EditorPageDTO `json:"pages"`
}

type SaveEditorRequest struct {
	Pages []EditorPageDTO `json:"pages"`
	DeletedPages  []int           `json:"deletedPages"`
}
type CreateUploadedBookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`

	CoverURL    string `json:"cover_url"`

	AccessType  string `json:"access_type"`

	AgeMin      int    `json:"age_min"`
	AgeMax      int    `json:"age_max"`

	Language    string `json:"language"`
	Category    string `json:"category"`
}