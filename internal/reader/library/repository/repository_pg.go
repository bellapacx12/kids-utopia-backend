package repository

import (
	"context"

	progressmodel "github.com/bellapacx/kids-utopia/internal/progress/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) GetContinueReading(
	ctx context.Context,
	childID string,
	limit int,
) ([]progressmodel.BookProgress, error) {

	query := `
		SELECT
			id,
			child_id,
			book_id,
			current_page,
			progress_percent,
			completed,
			last_read_at,
			created_at,
			updated_at
		FROM book_progress
		WHERE child_id = $1
		  AND completed = false
		ORDER BY last_read_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(
		ctx,
		query,
		childID,
		limit,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []progressmodel.BookProgress

	for rows.Next() {

		var p progressmodel.BookProgress

		err := rows.Scan(
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
			return nil, err
		}

		items = append(items, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}