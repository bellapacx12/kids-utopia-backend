package milestones

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/gamification/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}
func (r *Repo) HasAwarded(
	ctx context.Context,
	childID, milestoneID, bookID string,
) (bool, error) {

	var exists bool

	err := r.db.QueryRow(ctx, `
	SELECT EXISTS (
		SELECT 1
		FROM milestones_awarded
		WHERE child_id = $1
		  AND milestone_id = $2
		  AND (
				(book_id = $3)
				OR ($3 IS NULL AND book_id IS NULL)
		  )
	)
`, childID, milestoneID, nullIfEmpty(bookID)).Scan(&exists)
	return exists, err
}
func (r *Repo) Award(
	ctx context.Context,
	a *model.MilestoneAward,
) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO milestones_awarded (
			id,
			child_id,
			book_id,
			milestone_id,
			reward_type,
			reward_ref,
			event_id,
			created_at
		)
		VALUES (
			DEFAULT,
			$1,$2,$3,$4,$5,$6,NOW()
		)
		ON CONFLICT DO NOTHING
	`,
		a.ChildID,
		a.BookID,
		a.MilestoneID,
		a.RewardType,
		a.RewardRef,
		a.EventID,
	)

	return err
}
func (r *Repo) InsertEvent(
	ctx context.Context,
	e *model.MilestoneEvent,
) (string, error) {

	var id string

	err := r.db.QueryRow(ctx, `
		INSERT INTO milestone_events (
			child_id,
			book_id,
			type,
			title,
			meta,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,NOW())
		RETURNING id
	`,
		e.ChildID,
		e.BookID,
		e.Type,
		e.Title,
		e.Meta,
	).Scan(&id)

	return id, err
}
func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}