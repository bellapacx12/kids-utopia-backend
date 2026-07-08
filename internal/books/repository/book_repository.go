package repository

import (
	"context"
	"database/sql"

	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/model"
	"github.com/bellapacx/kids-utopia/pkg/database"
)

type BookRepository interface {
	Create(ctx context.Context, book *model.Book) error
	FindByID(ctx context.Context, id string) (*model.Book, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	ListBooks(ctx context.Context,) ([]model.Book, error)
    
	// ✅ ADD THIS
	GetBookByID(ctx context.Context, id string) (*model.Book, error)

	// ✅ ADD THIS
	GetBookPages(ctx context.Context, bookID string) ([]model.BookPage, error)

	// ✅ ADD THIS
	GetBookPreview(ctx context.Context, bookID string) ([]model.BookPage, error)
	CreateVariant(ctx context.Context, v *model.BookVariant) error 
	FindVariantByID(
	ctx context.Context,
	id string,
) (*model.BookVariant, error) 
ListVariantsByBookID(ctx context.Context, bookID string) ([]model.BookVariant, error)
GetPagesByVariantID(
    ctx context.Context,
    variantID string,
) ([]dto.EditorPageDTO, error) 
GetPagesByVariantIDD(
    ctx context.Context,
    variantID string,
) ([]model.BookPage, error)
}
func (r *bookRepository) ListBooks(
	ctx context.Context,
) ([]model.Book, error) {

	query := `
		SELECT 
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
			progress,
			created_at,
			updated_at
		FROM books
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []model.Book

	for rows.Next() {
		var b model.Book

		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.Author,
			&b.CoverURL,
			&b.Status,
			&b.AccessType,
			&b.AgeMin,
			&b.AgeMax,
			&b.Language,
			&b.Category,
			&b.PopularityScore,
			&b.Progress,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
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
func (r *bookRepository) CreateVariant(ctx context.Context, v *model.BookVariant) error {

	query := `
	INSERT INTO book_variants (
		id, book_id, language,
		file_url, title,
		status, progress,
		created_at, updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := database.DB.Exec(ctx, query,
		v.ID,
		v.BookID,
		v.Language,
		v.FileURL,
		v.Title,
		v.Status,
		v.Progress,
		v.CreatedAt,
		v.UpdatedAt,
	)

	return err
}
func (r *bookRepository) FindVariantByID(
	ctx context.Context,
	id string,
) (*model.BookVariant, error) {

	query := `
		SELECT
			id,
			book_id,
			language,
			title,
			file_url,
			status,
			progress,
			created_at,
			updated_at
		FROM book_variants
		WHERE id = $1
	`

	row := database.DB.QueryRow(ctx, query, id)

	var v model.BookVariant

	err := row.Scan(
		&v.ID,
		&v.BookID,
		&v.Language,
		&v.Title,
		&v.FileURL,
		&v.Status,
		&v.Progress,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &v, nil
}
func (r *bookRepository) ListVariantsByBookID(
	ctx context.Context,
	bookID string,
) ([]model.BookVariant, error) {

	query := `
	SELECT
		id,
		book_id,
		language,
		title,
		file_url,
		status,
		progress,
		created_at,
		updated_at
	FROM book_variants
	WHERE book_id = $1
	ORDER BY created_at ASC
	`

	rows, err := database.DB.Query(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []model.BookVariant

	for rows.Next() {
		var v model.BookVariant

		err := rows.Scan(
			&v.ID,
			&v.BookID,
			&v.Language,
			&v.Title,
			&v.FileURL,
			&v.Status,
			&v.Progress,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		variants = append(variants, v)
	}

	return variants, nil
}
func (r *bookRepository) GetPagesByVariantID(
    ctx context.Context,
    variantID string,
) ([]dto.EditorPageDTO, error) {

	rows, err := database.DB.Query(ctx, `
		SELECT
			page_number,
			content,
			image_key,
			image_url
		FROM book_pages
		WHERE variant_id = $1
		ORDER BY page_number
	`, variantID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []dto.EditorPageDTO

	for rows.Next() {
		var p dto.EditorPageDTO

		err := rows.Scan(
			&p.PageNumber,
			&p.Content,
			&p.ImageKey,
			&p.ImageURL,
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}
func (r *bookRepository) GetPagesByVariantIDD(
	ctx context.Context,
	variantID string,
) ([]model.BookPage, error) {

	rows, err := database.DB.Query(ctx, `
		SELECT
			page_number,
			content,
			image_url,
			audio_url
		FROM book_pages
		WHERE variant_id = $1
		ORDER BY page_number
	`, variantID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []model.BookPage

	for rows.Next() {
		var p model.BookPage
		var audioURL sql.NullString

		err := rows.Scan(
			&p.PageNumber,
			&p.Content,
			&p.ImageURL,
			&audioURL,
		)
		if err != nil {
			return nil, err
		}

		if audioURL.Valid {
			p.AudioURL = &audioURL.String
		}

		pages = append(pages, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}