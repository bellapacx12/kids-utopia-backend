package repository

import (
	"context"
	"fmt"

	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookPagesRepository interface {
	GetPages(ctx context.Context, bookID string) ([]dto.EditorPageDTO, error)
	SavePages(ctx context.Context, bookID string, pages []dto.EditorPageDTO) error
	UpdateCoverURL(ctx context.Context, bookID string, coverURL string) error
	UpdateProgress(
    ctx context.Context,
    id string,
    status string,
    progress int,
) error
 SavePagesByVariant(
	ctx context.Context,
	variantID string,
	pages []dto.EditorPageDTO,
) error
GetPagesByVariantID(
    ctx context.Context,
    variantID string,
) ([]dto.EditorPageDTO, error)
}

type bookPagesRepo struct {
	db *pgxpool.Pool
}


func NewBookPagesRepository(db *pgxpool.Pool) BookPagesRepository {
	return &bookPagesRepo{db: db}
}
func (r *bookPagesRepo) GetPages(
	ctx context.Context,
	bookID string,
) ([]dto.EditorPageDTO, error) {

	query := `
		SELECT page_number, content, image_key, image_url
		FROM book_pages
		WHERE book_id = $1
		ORDER BY page_number ASC
	`

	rows, err := r.db.Query(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []dto.EditorPageDTO

	for rows.Next() {
		var p dto.EditorPageDTO

		if err := rows.Scan(
			&p.PageNumber,
			&p.Content,
			&p.ImageKey,
			&p.ImageURL,
		); err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, rows.Err()
}
func (r *bookPagesRepo) SavePages(
	ctx context.Context,
	bookID string,
	pages []dto.EditorPageDTO,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// clear old version
	_, err = tx.Exec(ctx,
		`DELETE FROM book_pages WHERE book_id = $1`,
		bookID,
	)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO book_pages
		(id, book_id, page_number, content, image_key, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, p := range pages {
		_, err := tx.Exec(ctx,
			query,
			uuid.NewString(),
			bookID,
			p.PageNumber,
			p.Content,
			p.ImageKey,
			p.ImageURL,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
func (r *bookPagesRepo) UpdateCoverURL(ctx context.Context, bookID string, coverURL string) error {
	query := `
		UPDATE books
		SET cover_url = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, coverURL, bookID)
	if err != nil {
		return fmt.Errorf("update cover_url failed: %w", err)
	}

	return nil
}
func (r *bookPagesRepo) UpdateProgress(
    ctx context.Context,
    id string,
    status string,
    progress int,
) error {

    _, err := r.db.Exec(ctx, `
       UPDATE book_variants
SET status = $1,
    progress = $2,
    updated_at = NOW()
WHERE id = $3
    `, status, progress, id)

    return err
}
func (r *bookPagesRepo) SavePagesByVariant(
	ctx context.Context,
	variantID string,
	pages []dto.EditorPageDTO,
) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// =========================
	// clear old pages for variant
	// =========================
	_, err = tx.Exec(ctx,
		`DELETE FROM book_pages WHERE variant_id = $1`,
		variantID,
	)
	if err != nil {
		return err
	}

	// =========================
	// insert new pages
	// =========================
	query := `
		INSERT INTO book_pages
		(id, variant_id, page_number, content, image_key, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, p := range pages {
		_, err := tx.Exec(ctx,
			query,
			uuid.NewString(),
			variantID,
			p.PageNumber,
			p.Content,
			p.ImageKey,
			p.ImageURL,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
func (r *bookPagesRepo) GetPagesByVariantID(
    ctx context.Context,
    variantID string,
) ([]dto.EditorPageDTO, error) {

	rows, err := r.db.Query(ctx, `
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
