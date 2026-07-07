package repository

import (
	"context"
	"errors"
	"time"

	"github.com/bellapacx/kids-utopia/internal/reader_session/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) SessionRepository {
	return &repo{db: db}
}
func (r *repo) Create(ctx context.Context, s *model.ReadingSession) error {

	s.ID = uuid.NewString()
	now := time.Now()

	s.StartedAt = now
	s.CreatedAt = now
	s.UpdatedAt = now

	query := `
		INSERT INTO reading_sessions (
			id,
			user_id,
			child_id,
			book_id,
			started_at,
			start_page,
			end_page,
			completed,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		s.ID,
		s.UserID,
		s.ChildID,
		s.BookID,
		s.StartedAt,
		s.StartPage,
		s.EndPage,
		s.Completed,
		s.CreatedAt,
		s.UpdatedAt,
	)

	return err
}
func (r *repo) GetByID(ctx context.Context, id string) (*model.ReadingSession, error) {

	var s model.ReadingSession

	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, child_id, book_id,
		       started_at, ended_at,
		       duration_seconds,
		       start_page, end_page,
		       completed,
		       created_at, updated_at
		FROM reading_sessions
		WHERE id = $1
	`, id).Scan(
		&s.ID,
		&s.UserID,
		&s.ChildID,
		&s.BookID,
		&s.StartedAt,
		&s.EndedAt,
		&s.DurationSeconds,
		&s.StartPage,
		&s.EndPage,
		&s.Completed,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}
func (r *repo) GetActiveSession(
	ctx context.Context,
	userID, childID, bookID string,
) (*model.ReadingSession, error) {

	var s model.ReadingSession

	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, child_id, book_id,
		       started_at, ended_at,
		       start_page, end_page,
		       completed
		FROM reading_sessions
		WHERE user_id = $1
        AND child_id = $2
        AND book_id = $3
        AND completed = false
        ORDER BY started_at DESC
        LIMIT 1
	`, userID, childID, bookID).Scan(
		&s.ID,
		&s.UserID,
		&s.ChildID,
		&s.BookID,
		&s.StartedAt,
		&s.EndedAt,
		&s.StartPage,
		&s.EndPage,
		&s.Completed,
	)

	if err != nil {

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return nil, err
}

	return &s, nil
}
func (r *repo) Update(ctx context.Context, s *model.ReadingSession) error {

	query := `
		UPDATE reading_sessions
		SET end_page = $1,
		    completed = $2,
		    updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(
		ctx,
		query,
		s.EndPage,
		s.Completed,
		s.ID,
	)

	return err
}
func (r *repo) EndSession(
	ctx context.Context,
	s *model.ReadingSession,
) error {

	now := time.Now()

	query := `
		UPDATE reading_sessions
		SET end_page = $1,
		    ended_at = $2,
		    updated_at = $3,
		    completed = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		ctx,
		query,
		s.EndPage,
		now,
		now,
		s.Completed,
		s.ID,
	)

	return err
}
func (r *repo) GetTotalReadingTime(ctx context.Context, childID string) (int, error) {
	var total int

	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(duration_seconds), 0)
		FROM reading_sessions
		WHERE child_id = $1
	`, childID).Scan(&total)

	return total, err
}
func (r *repo) CountSessions(ctx context.Context, childID string) (int, error) {
	var count int

	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM reading_sessions
		WHERE child_id = $1
	`, childID).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}