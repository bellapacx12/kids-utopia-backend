package repository

import (
	"context"
	"time"

	"github.com/bellapacx/kids-utopia/internal/progress/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgressRepository interface {
	Create(ctx context.Context, p *model.BookProgress) error
	Get(ctx context.Context, childID, bookID string) (*model.BookProgress, error)
	Update(ctx context.Context, p *model.BookProgress) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewProgressRepository(db *pgxpool.Pool) ProgressRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, p *model.BookProgress) error {
	now := time.Now()

	_, err := r.db.Exec(ctx, `
		INSERT INTO book_progress
		(child_id, book_id, current_page, progress_percent, completed, last_read_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		p.ChildID,
		p.BookID,
		p.CurrentPage,
		p.ProgressPercent,
		p.Completed,
		now,
		now,
		now,
	)

	return err
}
func (r *repo) Get(ctx context.Context, childID, bookID string) (*model.BookProgress, error) {
	var p model.BookProgress

	err := r.db.QueryRow(ctx, `
		SELECT id, child_id, book_id, current_page, progress_percent, completed, last_read_at, created_at, updated_at
		FROM book_progress
		WHERE child_id=$1 AND book_id=$2
	`, childID, bookID).Scan(
		&p.ID,
		&p.ChildID,
		&p.BookID,
		&p.CurrentPage,
		&p.ProgressPercent,
		&p.Completed,
		&p.LastReadAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	return nil, err
}

	return &p, nil
}
func (r *repo) Update(ctx context.Context, p *model.BookProgress) error {
	now := time.Now()

	_, err := r.db.Exec(ctx, `
		UPDATE book_progress
		SET current_page=$1,
		    progress_percent=$2,
		    completed=$3,
		    last_read_at=$4,
		    updated_at=$5
		WHERE child_id=$6 AND book_id=$7
	`,
		p.CurrentPage,
		p.ProgressPercent,
		p.Completed,
		now,
		now,
		p.ChildID,
		p.BookID,
	)

	return err
}