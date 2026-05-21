package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/subscriptions/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, sub model.Subscription) error
	GetActiveByUser(ctx context.Context, userID string) (*model.Subscription, error)
}

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}
func (r *repo) Create(ctx context.Context, sub model.Subscription) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO subscriptions (
			id, user_id, plan, status, start_date, end_date, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())
	`,
		sub.ID,
		sub.UserID,
		sub.Plan,
		sub.Status,
		sub.StartDate,
		sub.EndDate,
	)

	return err
}
func (r *repo) GetActiveByUser(ctx context.Context, userID string) (*model.Subscription, error) {

	var sub model.Subscription

	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, plan, status, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		AND status = 'active'
		AND (end_date IS NULL OR end_date > NOW())
		LIMIT 1
	`, userID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.Plan,
		&sub.Status,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}