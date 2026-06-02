package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/analytics/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}
func (r *Repo) Insert(ctx context.Context, e model.Event) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO analytics_events
		(id, type, user_id, child_id, book_id, session_id, page, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (id) DO NOTHING
	`,
		e.EventID,
		e.Type,
		e.UserID,
		e.ChildID,
		e.BookID,
		e.SessionID,
		e.Page,
		e.Timestamp,
	)

	return err
}
func (r *Repo) GetChildStats(ctx context.Context, childID string) (*model.Stats, error) {

	var stats model.Stats

	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE type = 'SESSION_STARTED') AS total_sessions,
			COALESCE(SUM(page) FILTER (WHERE type = 'PROGRESS_UPDATED'), 0) AS total_pages
		FROM analytics_events
		WHERE child_id = $1
	`, childID).Scan(
		&stats.TotalSessions,
		&stats.TotalPages,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}