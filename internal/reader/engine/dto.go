package engine

type OpenRequest struct {
	ChildID string `json:"child_id"`
	BookID  string `json:"book_id"`
}

type UpdateRequest struct {
	ChildID string `json:"child_id"`
	BookID  string `json:"book_id"`
	Page    int    `json:"page"`
}

type CloseRequest struct {
	ChildID string `json:"child_id"`
	BookID  string `json:"book_id"`
	Page    int    `json:"page"`
}
