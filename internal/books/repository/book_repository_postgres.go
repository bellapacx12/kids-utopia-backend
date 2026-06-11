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
func (r *bookRepository) Create(
	ctx context.Context,
	b *model.Book,
) error {

	query := `
	INSERT INTO books (
		id,
		title,
		description,
		author,
		cover_url,
		status,
		access_type,
		age_min,
		age_max,
		language,
		category,
		popularity_score,
		created_at,
		updated_at
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,
		$7,$8,$9,$10,$11,$12,
		$13,$14
	)
	`

	_, err := database.DB.Exec(
		ctx,
		query,
		b.ID,
		b.Title,
		b.Description,
		b.Author,
		b.CoverURL,
		b.Status,

		b.AccessType,
		b.AgeMin,
		b.AgeMax,
		b.Language,
		b.Category,
		b.PopularityScore,

		b.CreatedAt,
		b.UpdatedAt,
	)

	return err
}
func (r *bookRepository) FindByID(ctx context.Context, id string) (*model.Book, error) {

	query := `
	SELECT 
		id,
		title,
		description,
		author,
		cover_url,
		status,
		progress,
		access_type,
		age_min,
		age_max,
		language,
		category,
		created_at,
		updated_at
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
		&b.Progress,
		&b.AccessType,
		&b.AgeMin,
		&b.AgeMax,
		&b.Language,
		&b.Category,
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

func (r *bookRepository) GetBookByID(
	ctx context.Context,
	id string,
) (*model.Book, error) {

	var b model.Book

	query := `
		SELECT id, title, author, cover_url, access_type, created_at
		FROM books
		WHERE id = $1
	`

	err := database.DB.QueryRow(ctx, query, id).Scan(
		&b.ID,
		&b.Title,
		&b.Author,
		&b.CoverURL,
		&b.AccessType,
		&b.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}
func (r *bookRepository) GetBookPages(
	ctx context.Context,
	bookID string,
) ([]model.BookPage, error) {

	query := `
		SELECT page_number, content, image_url
		FROM book_pages
		WHERE book_id = $1
		ORDER BY page_number ASC
	`

	rows, err := database.DB.Query(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []model.BookPage

	for rows.Next() {
		var p model.BookPage

		if err := rows.Scan(
			&p.PageNumber,
			&p.Content,
			&p.ImageURL,
		); err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}
func (r *bookRepository) GetBookPreview(
	ctx context.Context,
	bookID string,
) ([]model.BookPage, error) {

	query := `
		SELECT page_number, content, image_url
		FROM book_pages
		WHERE book_id = $1
		ORDER BY page_number ASC
		LIMIT 2
	`

	rows, err := database.DB.Query(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []model.BookPage

	for rows.Next() {
		var p model.BookPage

		if err := rows.Scan(
			&p.PageNumber,
			&p.Content,
			&p.ImageURL,
		); err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}
