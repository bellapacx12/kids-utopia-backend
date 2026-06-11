package service

import (
	"context"
	"sort"

	bookmarkrepo "github.com/bellapacx/kids-utopia/internal/bookmarks/repository"
	bookrepo "github.com/bellapacx/kids-utopia/internal/books/repository"
	childrepo "github.com/bellapacx/kids-utopia/internal/children/repository"
	progressrepo "github.com/bellapacx/kids-utopia/internal/progress/repository"
	"github.com/bellapacx/kids-utopia/internal/recommendation/dto"
)

type Service struct {
	bookRepo     bookrepo.BookRepository
	progressRepo progressrepo.ProgressRepository
	bookmarkRepo bookmarkrepo.Repository
	childRepo   childrepo.ChildRepository
}

func New(
	b bookrepo.BookRepository,
	p progressrepo.ProgressRepository,
	bo bookmarkrepo.Repository,
	c childrepo.ChildRepository,
) *Service {
	return &Service{
		bookRepo:     b,
		progressRepo: p,
		bookmarkRepo: bo,
		childRepo: c,
	}
}
func (s *Service) Recommend(
	ctx context.Context,
	childID string,
) ([]dto.Recommendation, error) {

	books, err := s.bookRepo.ListBooks(ctx)
	if err != nil {
		return nil, err
	}

	progress, _ := s.progressRepo.ListByChild(ctx, childID)
	bookmarks, _ := s.bookmarkRepo.ListByChild(ctx, childID)

	// 👉 get child profile (IMPORTANT CHANGE)
	child, err := s.childRepo.FindByID(ctx, childID)
	if err != nil {
		return nil, err
	}

	completed := map[string]bool{}
	for _, p := range progress {
		if p.Completed {
			completed[p.BookID] = true
		}
	}

	result := make([]dto.Recommendation, 0)

	for _, b := range books {

		// skip completed
		if completed[b.ID] {
			continue
		}

		score := 0

		// AGE MATCH (backend controlled)
		if child.Age >= b.AgeMin && child.Age <= b.AgeMax {
			score += 40
		}

		// LANGUAGE MATCH (backend controlled)
		if b.Language == child.Language {
			score += 20
		}

		// CATEGORY BOOST
		if b.Category != "" {
			score += 10
		}

		// BOOKMARK BOOST
		for _, bk := range bookmarks {
			if bk.BookID == b.ID {
				score += 30
				break
			}
		}

		// POPULARITY BOOST
		score += b.PopularityScore / 10

		result = append(result, dto.Recommendation{
			BookID:   b.ID,
			Title:    b.Title,
			CoverURL: b.CoverURL,
			Score:    score,
			Reason:   "Recommended for you",
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result, nil
}