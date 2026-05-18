package service

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/users/dto"
	"github.com/bellapacx/kids-utopia/internal/users/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}
func (s *UserService) GetMe(ctx context.Context, userID string) (*dto.UserResponse, error) {

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Phone:      user.Phone,
		Name:       user.Name,
		AvatarURL:  user.AvatarURL,
		Role:       user.Role,
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
	}, nil
}
func (s *UserService) UpdateMe(ctx context.Context, userID string, req dto.UpdateProfileRequest) error {

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		user.Name = req.Name
	}

	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	return s.repo.Update(ctx, user)
}