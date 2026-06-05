package themes

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Unlock(ctx context.Context, childID, themeID string) error
	GetUnlocked(ctx context.Context, childID string) ([]string, error)
	GetUnlockedThemes(ctx context.Context, childID string) ([]string, error)
}
type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) Unlock(ctx context.Context, childID, themeID string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO child_themes (child_id, theme_id)
		VALUES ($1, $2)
		ON CONFLICT (child_id, theme_id) DO NOTHING
	`, childID, themeID)

	return err
}

func (r *repo) GetUnlocked(ctx context.Context, childID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT theme_id
		FROM child_themes
		WHERE child_id = $1
	`, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}
func (r *repo) GetUnlockedThemes(ctx context.Context, childID string) ([]string, error) {

	rows, err := r.db.Query(ctx, `
		SELECT theme_id
		FROM child_themes
		WHERE child_id = $1
	`, childID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var themeID string
		if err := rows.Scan(&themeID); err != nil {
			return nil, err
		}
		result = append(result, themeID)
	}

	return result, nil
}