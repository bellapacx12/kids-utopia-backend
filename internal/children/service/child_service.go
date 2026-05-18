package service

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/children/dto"
	"github.com/bellapacx/kids-utopia/internal/children/model"
	"github.com/bellapacx/kids-utopia/internal/children/repository"
)

type ChildService struct {
	repo repository.ChildRepository
}

func NewChildService(r repository.ChildRepository) *ChildService {
	return &ChildService{repo: r}
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