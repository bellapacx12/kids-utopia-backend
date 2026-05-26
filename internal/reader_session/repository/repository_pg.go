package repository

import (
	"context"
	"time"

	"github.com/bellapacx/kids-utopia/internal/reader_session/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
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
			duration_seconds,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
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
		s.DurationSeconds,
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
		  AND ended_at IS NULL
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

	duration := int(now.Sub(s.StartedAt).Seconds())

	query := `
		UPDATE reading_sessions
		SET end_page = $1,
		    completed = true,
		    ended_at = $2,
		    duration_seconds = $3,
		    updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		ctx,
		query,
		s.EndPage,
		now,
		duration,
		now,
		s.ID,
	)

	return err
}