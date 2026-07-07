package service

import (
	"context"
	"errors"
	"log"
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

	log.Println("========== UpdateProgress ==========")
	log.Printf("INPUT childID=%s bookID=%s page=%d totalPages=%d",
		childID,
		bookID,
		page,
		totalPages,
	)

	progress, err := s.repo.Get(
		ctx,
		childID,
		bookID,
	)

	// =========================
	// CREATE NEW PROGRESS
	// =========================

	if err != nil {

		log.Printf("repo.Get returned error: %v", err)

		if errors.Is(err, repository.ErrNotFound) {

			log.Println("No existing progress found. Creating new progress...")

			percent := calculateProgress(
				page,
				totalPages,
			)

			log.Printf("Calculated progress percent=%d", percent)

			p := &model.BookProgress{
				ChildID:         childID,
				BookID:          bookID,
				CurrentPage:     page,
				ProgressPercent: percent,
				Completed:       percent >= 100,
				LastReadAt:      time.Now(),
			}

			log.Printf("Creating progress: %+v", p)

			err := s.repo.Create(ctx, p)
			if err != nil {
				log.Printf("repo.Create failed: %v", err)
				return err
			}

			log.Println("Progress created successfully.")
			return nil
		}

		log.Printf("repo.Get failed: %v", err)
		return err
	}

	log.Printf("Existing progress:")
	log.Printf("  CurrentPage     = %d", progress.CurrentPage)
	log.Printf("  ProgressPercent = %d", progress.ProgressPercent)
	log.Printf("  Completed       = %v", progress.Completed)

	if progress.CurrentPage == page {
		log.Printf(
			"Skipping update because current page (%d) == incoming page (%d)",
			progress.CurrentPage,
			page,
		)
		return nil
	}

	log.Printf(
		"Updating page from %d -> %d",
		progress.CurrentPage,
		page,
	)

	progress.CurrentPage = page

	progress.ProgressPercent = calculateProgress(
		page,
		totalPages,
	)

	log.Printf(
		"New ProgressPercent=%d",
		progress.ProgressPercent,
	)

	progress.Completed = progress.ProgressPercent >= 100

	log.Printf(
		"Completed=%v",
		progress.Completed,
	)

	progress.LastReadAt = time.Now()

	log.Printf("Saving progress: %+v", progress)

	err = s.repo.Update(ctx, progress)
	if err != nil {
		log.Printf("repo.Update failed: %v", err)
		return err
	}

	log.Println("Progress updated successfully.")
	log.Println("====================================")

	return nil
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
func (s *ProgressService) Exists(
	ctx context.Context,
	childID, bookID string,
) (bool, error) {
	return s.repo.Exists(ctx, childID, bookID)
}