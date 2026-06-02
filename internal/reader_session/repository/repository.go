package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/reader_session/model"
)

type SessionRepository interface {
	Create(ctx context.Context, s *model.ReadingSession) error

	GetByID(ctx context.Context, id string) (*model.ReadingSession, error)

	GetActiveSession(
		ctx context.Context,
		userID, childID, bookID string,
	) (*model.ReadingSession, error)

	Update(ctx context.Context, s *model.ReadingSession) error

	EndSession(ctx context.Context, s *model.ReadingSession) error
	GetTotalReadingTime(ctx context.Context, childID string) (int, error) 
}