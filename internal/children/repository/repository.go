package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/children/model"
)

type ChildRepository interface {
	Create(ctx context.Context, child *model.Child) error
	FindByParentID(ctx context.Context, parentID string) ([]model.Child, error)
	FindByID(ctx context.Context, id string) (*model.Child, error)
}