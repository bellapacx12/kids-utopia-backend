package dto

type CreateChildRequest struct {
	Name      string  `json:"name" binding:"required"`
	AvatarURL *string `json:"avatar_url"`
	Age       *int    `json:"age"`
	Language  string  `json:"language"`
}