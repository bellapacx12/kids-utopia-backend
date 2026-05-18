package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/users/model"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}