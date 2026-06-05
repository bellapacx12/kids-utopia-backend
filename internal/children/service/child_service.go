package service

import (
	"context"
	"log"

	"github.com/bellapacx/kids-utopia/internal/children/dto"
	"github.com/bellapacx/kids-utopia/internal/children/model"
	"github.com/bellapacx/kids-utopia/internal/children/repository"
	GamificationService "github.com/bellapacx/kids-utopia/internal/gamification/service"
)

type ChildService struct {
	repo repository.ChildRepository
	gamification *GamificationService.Service
}

func NewChildService(r repository.ChildRepository, gamification *GamificationService.Service,) *ChildService {
	return &ChildService{repo: r, gamification: gamification}
}
func (s *ChildService) Create(ctx context.Context, parentID string, req dto.CreateChildRequest) error {

	child := &model.Child{
		ParentID: parentID,
		Name:     req.Name,
		AvatarURL: req.AvatarURL,
		Age:      req.Age,
		Language: req.Language,
	}

	if child.Language == "" {
		child.Language = "en"
	}

	return s.repo.Create(ctx, child)
}
func (s *ChildService) GetByParent(ctx context.Context, parentID string) ([]model.Child, error) {
	return s.repo.FindByParentID(ctx, parentID)
}
func (s *ChildService) FindByID(ctx context.Context, id string) (*model.Child, *GamificationService.Snapshot, error) {

	child, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	snap, err := s.gamification.GetSnapshot(ctx, id)
	if err != nil {
		log.Println("gamification error:", err)
		return child, nil, nil // never block child
	}

	return child, snap, nil
}