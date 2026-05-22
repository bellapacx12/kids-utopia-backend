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

	// ✅ ADD THIS
	GetBookByID(ctx context.Context, id string) (*model.Book, error)

	// ✅ ADD THIS
	GetBookPages(ctx context.Context, bookID string) ([]model.BookPage, error)

	// ✅ ADD THIS
	GetBookPreview(ctx context.Context, bookID string) ([]model.BookPage, error)
}
func (r *bookRepository) ListBooks(
	ctx context.Context,
	limit int,
	offset int,
) ([]model.Book, int, error) {

	query := `
		SELECT id, title, author, cover_url, status, created_at, access_type
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
			&b.AccessType,
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
func (r *bookRepository) UpdateAccessType(
	ctx context.Context,
	bookID string,
	accessType string,
) error {

	query := `
		UPDATE books
		SET access_type = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := database.DB.Exec(ctx, query, accessType, bookID)
	if err != nil {
		return err
	}

	return nil
}