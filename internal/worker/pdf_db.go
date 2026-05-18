package worker

import (
	"context"

	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/google/uuid"
)

func savePages(ctx context.Context, bookID string, pages []string, imageURLs []string) error {

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, content := range pages {

		var imageURL string
		if i < len(imageURLs) {
			imageURL = imageURLs[i]
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO book_pages (
				id,
				book_id,
				page_number,
				content,
				image_url,
				created_at
			)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW())
		`,
			bookID,
			i+1,
			content,
			imageURL,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
func saveDraftPages(
	ctx context.Context,
	bookID string,
	textPages []string,
	imageURLs []string,
) error {

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i := range imageURLs {

		var text string
		if i < len(textPages) {
			text = textPages[i]
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO book_raw_pages (
				id,
				book_id,
				page_number,
				image_url,
				text
			) VALUES ($1,$2,$3,$4,$5)
		`,
			uuid.NewString(),
			bookID,
			i+1,
			imageURLs[i],
			text,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}