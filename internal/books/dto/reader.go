package dto

import (
	"github.com/bellapacx/kids-utopia/internal/books/model"
)
type ReaderVariant struct {
    ID       string               `json:"id"`
    Title    string          `json:"title"`
    Language string               `json:"language"`
    Pages    []model.BookPage `json:"pages"`
}