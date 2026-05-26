package dto

type StartSessionRequest struct {
	ChildID string `json:"child_id"`
	BookID  string `json:"book_id"`
	Page    int    `json:"page"`
}

type UpdateSessionRequest struct {
	SessionID string `json:"session_id"`
	Page      int    `json:"page"`
}

type EndSessionRequest struct {
	SessionID string `json:"session_id"`
	Page      int    `json:"page"`
}