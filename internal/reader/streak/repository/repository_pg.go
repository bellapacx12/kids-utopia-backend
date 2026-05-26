package repository

import (
	"context"
	"time"

	"github.com/bellapacx/kids-utopia/internal/reader/streak/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) StreakRepository {
	return &repo{db: db}
}
func (r *repo) Get(ctx context.Context, childID string) (*model.ReadingStreak, error) {

	var s model.ReadingStreak

	err := r.db.QueryRow(ctx, `
		SELECT child_id, current_streak, longest_streak, last_read_date, updated_at
		FROM reading_streaks
		WHERE child_id = $1
	`, childID).Scan(
		&s.ChildID,
		&s.CurrentStreak,
		&s.LongestStreak,
		&s.LastReadDate,
		&s.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
func (r *repo) Create(ctx context.Context, s *model.ReadingStreak) error {

	now := time.Now()

	_, err := r.db.Exec(ctx, `
		INSERT INTO reading_streaks
		(child_id, current_streak, longest_streak, last_read_date, updated_at)
		VALUES ($1,$2,$3,$4,$5)
	`,
		s.ChildID,
		s.CurrentStreak,
		s.LongestStreak,
		s.LastReadDate,
		now,
	)

	return err
}
func (r *repo) Update(ctx context.Context, s *model.ReadingStreak) error {

	now := time.Now()

	_, err := r.db.Exec(ctx, `
		UPDATE reading_streaks
		SET current_streak = $1,
		    longest_streak = $2,
		    last_read_date = $3,
		    updated_at = $4
		WHERE child_id = $5
	`,
		s.CurrentStreak,
		s.LongestStreak,
		s.LastReadDate,
		now,
		s.ChildID,
	)

	return err
}