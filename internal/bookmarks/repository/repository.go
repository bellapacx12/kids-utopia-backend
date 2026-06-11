package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/bookmarks/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, b *model.Bookmark) error

	Delete(
		ctx context.Context,
		childID string,
		bookID string,
		page int,
	) error

	ListByBook(
		ctx context.Context,
		childID string,
		bookID string,
	) ([]model.Bookmark, error)
	ListByChild(
    ctx context.Context,
    childID string,
) ([]model.Bookmark, error)
ListDetailedByChild(
		ctx context.Context,
		childID string,
	) ([]model.BookmarkDetail, error)
}
type Repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}
func (r *Repo) Create(
	ctx context.Context,
	b *model.Bookmark,
) error {

	query := `
		INSERT INTO bookmarks (
			child_id,
			book_id,
			page
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		b.ChildID,
		b.BookID,
		b.Page,
	).Scan(
		&b.ID,
		&b.CreatedAt,
	)
}

func (r *Repo) Delete(
	ctx context.Context,
	childID string,
	bookID string,
	page int,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM bookmarks
		WHERE child_id = $1
		  AND book_id = $2
		  AND page = $3
		`,
		childID,
		bookID,
		page,
	)

	return err
}

func (r *Repo) ListByBook(
	ctx context.Context,
	childID string,
	bookID string,
) ([]model.Bookmark, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			child_id,
			book_id,
			page,
			created_at
		FROM bookmarks
		WHERE child_id = $1
		  AND book_id = $2
		ORDER BY page ASC
		`,
		childID,
		bookID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Bookmark

	for rows.Next() {

		var b model.Bookmark

		err := rows.Scan(
			&b.ID,
			&b.ChildID,
			&b.BookID,
			&b.Page,
			&b.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, b)
	}

	return result, rows.Err()
}
func (r *Repo) ListByChild(
	ctx context.Context,
	childID string,
) ([]model.Bookmark, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			child_id,
			book_id,
			page,
			created_at
		FROM bookmarks
		WHERE child_id = $1
		ORDER BY created_at DESC
		`,
		childID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Bookmark

	for rows.Next() {
		var b model.Bookmark

		if err := rows.Scan(
			&b.ID,
			&b.ChildID,
			&b.BookID,
			&b.Page,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, b)
	}

	return result, rows.Err()
}
func (r *Repo) ListDetailedByChild(
	ctx context.Context,
	childID string,
) ([]model.BookmarkDetail, error) {

	rows, err := r.db.Query(
		ctx,
		`
	SELECT
		b.id,
		b.book_id,
		bk.title,
		bk.cover_url,
		b.page,
		bp.content,
		bp.image_url,
		b.created_at
	FROM bookmarks b
	INNER JOIN books bk
		ON bk.id = b.book_id
	INNER JOIN book_pages bp
		ON bp.book_id = b.book_id
		AND bp.page_number = b.page
	WHERE b.child_id = $1
	ORDER BY b.created_at DESC
`,
		childID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.BookmarkDetail

	for rows.Next() {

		var b model.BookmarkDetail

		err := rows.Scan(
	&b.ID,
	&b.BookID,
	&b.BookTitle,
	&b.CoverImage,
	&b.Page,
	&b.Content,
	&b.ImageURL,
	&b.CreatedAt,
)

		if err != nil {
			return nil, err
		}

		result = append(result, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}