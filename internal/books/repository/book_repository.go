package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/books/model"
)

type BookRepository interface {
	Create(ctx context.Context, book *model.Book) error
	FindByID(ctx context.Context, id string) (*model.Book, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}