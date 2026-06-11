package service

import (
	"context"
	"fmt"
	"log"

	accesssvc "github.com/bellapacx/kids-utopia/internal/access/service"
	"github.com/bellapacx/kids-utopia/internal/bookmarks/model"
	"github.com/bellapacx/kids-utopia/internal/bookmarks/repository"
	bookmodel "github.com/bellapacx/kids-utopia/internal/books/model"
)

type Service struct {
	repo repository.Repository
	access *accesssvc.Service
}

func New(repo repository.Repository,  access *accesssvc.Service) *Service {
	return &Service{repo: repo, access: access}
}
func (s *Service) Create(
	ctx context.Context,
	b *model.Bookmark,
) error {

	// TODO: business rule checks can go here later
	// e.g. duplicate prevention, page validation

	return s.repo.Create(ctx, b)
}
func (s *Service) Delete(
	ctx context.Context,
	childID, bookID string,
	page int,
) error {

	return s.repo.Delete(ctx, childID, bookID, page)
}
func (s *Service) ListByBook(
	ctx context.Context,
	childID, bookID string,
) ([]model.Bookmark, error) {

	return s.repo.ListByBook(ctx, childID, bookID)
}
func (s *Service) ListByChild(
	ctx context.Context,
	childID string,
) ([]model.Bookmark, error) {

	return s.repo.ListByChild(ctx, childID)
}
func (s *Service) ListDetailedByChild(
    ctx context.Context,
    userID string,
    childID string,
) ([]model.BookmarkDetail, error) {

    // We need at least one book to validate access model,
    // so we should fetch any bookmarked book first or skip strict check

    bookmarks, err := s.repo.ListDetailedByChild(ctx, childID)
    if err != nil {
		log.Println("eroor here")
        return nil, err
    }

    if len(bookmarks) == 0 {
        return bookmarks, nil
    }

    // access check using first bookmark's book
    book := &bookmodel.Book{
        ID: bookmarks[0].BookID,
    }

    allowed, preview, err := s.access.CanAccessBook(ctx, userID, book)
if err != nil {
    return nil, err
}

if !allowed && !preview {
    return nil, fmt.Errorf("access denied")
}

    return bookmarks, nil
}