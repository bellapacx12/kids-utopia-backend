package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/books/model"
	"github.com/bellapacx/kids-utopia/pkg/database"
)

type BookRepository interface {
	Create(ctx context.Context, book *model.Book) error
	FindByID(ctx context.Context, id string) (*model.Book, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	ListBooks(ctx context.Context, limit int, offset int) ([]model.Book, int, error)
}
func (r *bookRepository) ListBooks(
	ctx context.Context,
	limit int,
	offset int,
) ([]model.Book, int, error) {

	query := `
		SELECT id, title, author, cover_url, status, created_at
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := database.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []model.Book

	for rows.Next() {
		var b model.Book

		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.CoverURL,
			&b.Status,
			&b.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		books = append(books, b)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	// =========================
	// COUNT QUERY
	// =========================

	var total int

	countQuery := `
		SELECT COUNT(*)
		FROM books
	`

	err = database.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}