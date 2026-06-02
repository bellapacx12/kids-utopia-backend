package service

import (
	"context"
	"time"

	"github.com/bellapacx/kids-utopia/internal/reader_session/model"
	"github.com/bellapacx/kids-utopia/internal/reader_session/repository"
)
type Service struct {
	repo repository.SessionRepository
	
}
func New(r repository.SessionRepository) *Service {
	return &Service{repo: r}
}
func (s *Service) StartSession(
	ctx context.Context,
	userID, childID, bookID string,
	page int,
) (*model.ReadingSession, error) {

	now := time.Now()

	session := &model.ReadingSession{
		UserID:    userID,
		ChildID:   childID,
		BookID:    bookID,
		StartPage: page,
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
		Completed:  false,
	}

	err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
func (s *Service) UpdateSession(
	ctx context.Context,
	sessionID string,
	page int,
) error {

	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	session.EndPage = &page
	session.UpdatedAt = time.Now()

	return s.repo.Update(ctx, session)
}
func (s *Service) EndSession(
	ctx context.Context,
	sessionID string,
	page int,
) error {

	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	now := time.Now()

	session.EndPage = &page
	session.EndedAt = &now
	session.Completed = true
	session.UpdatedAt = now
	session.DurationSeconds = int(now.Sub(session.StartedAt).Seconds())

	return s.repo.EndSession(ctx, session)
}
func (s *Service) GetSession(
	ctx context.Context,
	id string,
) (*model.ReadingSession, error) {

	return s.repo.GetByID(ctx, id)
}
func (s *Service) GetActiveSession(
	ctx context.Context,
	userID string,
	childID string,
	bookID string,
) (*model.ReadingSession, error) {
    
	return s.repo.GetActiveSession(
		ctx,
		userID,
		childID,
		bookID,
	)
}
func (s *Service) GetOrCreateActiveSession(
	ctx context.Context,
	userID, childID, bookID string,
	startPage int,
) (*model.ReadingSession, error) {

	// 1. Try get active session
	session, err := s.repo.GetActiveSession(ctx, userID, childID, bookID)
	if err != nil {
		return nil, err
	}

	// 2. If exists → reuse
	if session != nil {
		return session, nil
	}

	// 3. Otherwise create new
	return s.StartSession(ctx, userID, childID, bookID, startPage)
}