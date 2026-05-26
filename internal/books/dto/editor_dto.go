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
	Pages  []EditorPageDTO `json:"pages"`
}

type SaveEditorRequest struct {
	Pages []EditorPageDTO `json:"pages"`
	DeletedPages  []int           `json:"deletedPages"`
}
