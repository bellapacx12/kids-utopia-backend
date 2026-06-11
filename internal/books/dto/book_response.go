package dto

import "github.com/bellapacx/kids-utopia/internal/books/model"

type BookResponse struct {
	ID          string `json:"id"`
	Title       []string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Status      string `json:"status"`
}
type CreateBookResponse struct {
	Book    *model.Book
	Variant *model.BookVariant
}
type CreateVariantRequest struct {
	BookID   string
	Language string
	Title    string
	FileURL  string
}
type BookWithVariants struct {
	Book     Book           `json:"book"`
	Variants []model.BookVariant  `json:"variants"`
}
type Book struct {
	ID        string   `json:"id"`
	CoverURL  string   `json:"cover_url"`
	Title []string `json:"title"`
}