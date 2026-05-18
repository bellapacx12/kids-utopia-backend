package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bellapacx/kids-utopia/internal/children/model"
)

type childRepo struct {
	db *pgxpool.Pool
}

func NewChildRepository(db *pgxpool.Pool) ChildRepository {
	return &childRepo{db: db}
}
func (r *childRepo) Create(ctx context.Context, c *model.Child) error {

	query := `
		INSERT INTO children (parent_id, name, avatar_url, age, language)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		c.ParentID,
		c.Name,
		c.AvatarURL,
		c.Age,
		c.Language,
	)

	return err
}
func (r *childRepo) FindByParentID(ctx context.Context, parentID string) ([]model.Child, error) {

	query := `
		SELECT id, parent_id, name, avatar_url, age, language, created_at, updated_at
		FROM children
		WHERE parent_id=$1
	`

	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []model.Child

	for rows.Next() {
		var c model.Child

		err := rows.Scan(
			&c.ID,
			&c.ParentID,
			&c.Name,
			&c.AvatarURL,
			&c.Age,
			&c.Language,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		children = append(children, c)
	}

	return children, nil
}
func (r *childRepo) FindByID(ctx context.Context, id string) (*model.Child, error) {

	query := `
		SELECT id, parent_id, name, avatar_url, age, language, created_at, updated_at
		FROM children
		WHERE id=$1
	`

	var c model.Child

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.ParentID,
		&c.Name,
		&c.AvatarURL,
		&c.Age,
		&c.Language,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}