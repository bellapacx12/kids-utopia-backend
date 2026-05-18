package dto

type CreateBookRequest struct {
	Title      string `json:"title" binding:"required"`
	Description string `json:"description"`
	Author     string `json:"author"`
}