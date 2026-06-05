package service

import (
	"context"
	"errors"
	"time"

	"github.com/bellapacx/kids-utopia/internal/progress/model"
	"github.com/bellapacx/kids-utopia/internal/progress/repository"
)

type ProgressService struct {
	repo repository.ProgressRepository
}

func NewProgressService(r repository.ProgressRepository) *ProgressService {
	return &ProgressService{repo: r}
}
func (s *ProgressService) UpdateProgress(
	ctx context.Context,
	childID string,
	bookID string,
	page int,
	totalPages int,
) error {

	progress, err := s.repo.Get(
		ctx,
		childID,
		bookID,
	)

	// =========================
	// CREATE NEW PROGRESS
	// =========================

	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {

			percent := calculateProgress(
				page,
				totalPages,
			)

			p := &model.BookProgress{
				ChildID:         childID,
				BookID:          bookID,
				CurrentPage:     page,
				ProgressPercent: percent,
				Completed:       percent >= 100,
				LastReadAt:      time.Now(),
			}

			return s.repo.Create(ctx, p)
		}

		return err
	}
	
if progress.CurrentPage == page {
	return nil
}


	progress.CurrentPage = page

	progress.ProgressPercent = calculateProgress(
		page,
		totalPages,
	)

	progress.Completed = progress.ProgressPercent >= 100

	progress.LastReadAt = time.Now()

	return s.repo.Update(ctx, progress)
}

func calculateProgress(
	page int,
	totalPages int,
) int {

	if totalPages <= 0 {
		return 0
	}

	percent := (page * 100) / totalPages

	if percent > 100 {
		return 100
	}

	return percent
}
func (s *ProgressService) GetProgress(
	ctx context.Context,
	childID, bookID string,
) (*model.BookProgress, error) {
	return s.repo.Get(ctx, childID, bookID)
}
func (s *ProgressService) CreateProgress(
	ctx context.Context,
	childID string,
	bookID string,
	page int,
) (*model.BookProgress, error) {

	p := &model.BookProgress{
		ChildID:         childID,
		BookID:          bookID,
		CurrentPage:     page,
		ProgressPercent: 0,
		Completed:       false,
		LastReadAt:      time.Now(),
	}

	err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
func (s *ProgressService) ListByChild(
	ctx context.Context,
	childID string,
) ([]model.BookProgress, error) {
	return s.repo.ListByChild(ctx, childID)
}