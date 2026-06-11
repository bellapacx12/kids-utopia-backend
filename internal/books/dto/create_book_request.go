package dto

type CreateBookRequest struct {
	Title      string `json:"title" binding:"required"`
	Description string `json:"description"`
	Author     string `json:"author"`
}
type UploadVariantRequest struct {
	Language string `form:"language" binding:"required"`
}
type CreateFirstVariantRequest struct {
	BookID      string
	Title       string
	Author      string
	Description string

	Language    string
	AccessType  string
	Category    string

	AgeMin      int
	AgeMax      int
}