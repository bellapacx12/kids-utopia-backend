package service

import (
	"context"
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

func (s *ProgressService) UpdateProgress(ctx context.Context, childID, bookID string, page int) error {

	progress, err := s.repo.Get(ctx, childID, bookID)
	if err != nil {
		// create new progress
		p := &model.BookProgress{
			ChildID:         childID,
			BookID:          bookID,
			CurrentPage:     page,
			ProgressPercent: calculate(page),
			Completed:       false,
			LastReadAt:      time.Now(),
		}

		if err := s.repo.Create(ctx, p); err != nil {
			return err
		}
		return nil
	}

	progress.CurrentPage = page
	progress.ProgressPercent = calculate(page)
	progress.LastReadAt = time.Now()

	if progress.ProgressPercent >= 100 {
		progress.Completed = true
	}

	return s.repo.Update(ctx, progress)
}

func (s *ProgressService) GetProgress(ctx context.Context, childID, bookID string) (*model.BookProgress, error) {
	return s.repo.Get(ctx, childID, bookID)
}

func calculate(page int) int {
	// placeholder logic (you will improve later using total pages)
	if page >= 10 {
		return 100
	}
	return page * 10
}