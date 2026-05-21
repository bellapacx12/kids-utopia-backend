package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/books/model"
	"github.com/bellapacx/kids-utopia/pkg/database"
)

type bookRepository struct{}

func NewBookRepository() BookRepository {
	return &bookRepository{}
}
func (r *bookRepository) Create(ctx context.Context, b *model.Book) error {

	query := `
	INSERT INTO books (id, title, description, author, cover_url, status, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := database.DB.Exec(ctx, query,
		b.ID,
		b.Title,
		b.Description,
		b.Author,
		b.CoverURL,
		b.Status,
		b.CreatedAt,
		b.UpdatedAt,
	)

	return err
}
func (r *bookRepository) FindByID(ctx context.Context, id string) (*model.Book, error) {

	query := `
	SELECT id, title, description, author, cover_url, status, created_at, updated_at
	FROM books
	WHERE id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var b model.Book

	err := row.Scan(
		&b.ID,
		&b.Title,
		&b.Description,
		&b.Author,
		&b.CoverURL,
		&b.Status,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}
func (r *bookRepository) UpdateStatus(ctx context.Context, id string, status string) error {

	query := `UPDATE books SET status=$1, updated_at=NOW() WHERE id=$2`

	_, err := database.DB.Exec(ctx, query, status, id)

	return err
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