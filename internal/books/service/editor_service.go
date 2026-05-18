package service

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/pkg/storage"
)

type EditorService struct {
	bookRepo  repository.BookRepository
	pageRepo  repository.BookPagesRepository
	storage   storage.Storage
}
func (s *EditorService) GetEditor(
	ctx context.Context,
	bookID string,
) (*dto.EditorResponse, error) {

	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	pages, err := s.pageRepo.GetPages(ctx, bookID)
	if err != nil {
		return nil, err
	}

	// convert image_key → presigned URL
	for i := range pages {
		if pages[i].ImageKey != "" {
			pages[i].ImageURL = s.storage.GetPublicURL(
	pages[i].ImageKey,
)
		}
	}

	return &dto.EditorResponse{
		BookID: book.ID,
		Status: book.Status,
		Pages:  pages,
	}, nil
}
func (s *EditorService) SaveEditor(
	ctx context.Context,
	bookID string,
	req dto.SaveEditorRequest,
) error {

	return s.pageRepo.SavePages(ctx, bookID, req.Pages)
}
func NewEditorService(
	bookRepo repository.BookRepository,
	pageRepo repository.BookPagesRepository,
	storage storage.Storage,
) *EditorService {

	return &EditorService{
		bookRepo: bookRepo,
		pageRepo: pageRepo,
		storage:  storage,
	}
}