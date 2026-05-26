package library

import (
	"context"

	bookservice "github.com/bellapacx/kids-utopia/internal/books/service"
	"github.com/bellapacx/kids-utopia/internal/reader/library/repository"
)

type Service struct {
	repo        repository.Repository
	bookService *bookservice.BookService
}

func New(
	repo repository.Repository,
	bookService *bookservice.BookService,
) *Service {
	return &Service{
		repo:        repo,
		bookService: bookService,
	}
}

func (s *Service) GetContinueReading(
	ctx context.Context,
	childID string,
) (*ContinueReadingResponse, error) {

	progressItems, err := s.repo.GetContinueReading(
		ctx,
		childID,
		20,
	)

	if err != nil {
		return nil, err
	}

	response := &ContinueReadingResponse{
		Items: []ContinueReadingItem{},
	}

	for _, item := range progressItems {

		book, err := s.bookService.GetBookMeta(
			ctx,
			item.BookID,
		)

		if err != nil {
			continue
		}

		response.Items = append(
			response.Items,
			ContinueReadingItem{
				BookID:          item.BookID,
				Title:           book.Title,
				CoverURL:        book.CoverURL,
				CurrentPage:     item.CurrentPage,
				ProgressPercent: item.ProgressPercent,
				LastReadAt:      item.LastReadAt,
			},
		)
	}

	return response, nil
}