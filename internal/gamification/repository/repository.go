package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bellapacx/kids-utopia/internal/gamification/model"
	"github.com/google/uuid"

	"github.com/bellapacx/kids-utopia/internal/gamification/rules"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetXP(ctx context.Context, childID string) (*model.ChildXP, error)
   
	
	UpsertXP(
		ctx context.Context,
		childID string,
		amount int,
	) error
    GetState(
		ctx context.Context,
		childID string,
		
	) (rules.State, error)
	InsertTransaction(
		ctx context.Context,
		tx *model.XPTransaction,
	) error
	 GetBadges(
	ctx context.Context,
	childID string,
	
) ([]model.MilestoneAward, error)
}
type Repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetXP(
	ctx context.Context,
	childID string,
) (*model.ChildXP, error) {

	var xp model.ChildXP

	err := r.db.QueryRow(ctx, `
		SELECT
			child_id,
			total_xp,
			level,
			created_at,
			updated_at
		FROM child_xp
		WHERE child_id = $1
	`, childID).Scan(
		&xp.ChildID,
		&xp.TotalXP,
		&xp.Level,
		&xp.CreatedAt,
		&xp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &xp, nil
}

func (r *Repo) InsertTransaction(
	ctx context.Context,
	tx *model.XPTransaction,
) error {

	if tx.ID == "" {
		tx.ID = uuid.NewString()
	}

	_, err := r.db.Exec(ctx, `
	INSERT INTO xp_transactions (
		id,
		child_id,
		source,
		source_id,
		xp_amount,
		created_at
	)
	VALUES ($1,$2,$3,$4,$5,NOW())
	ON CONFLICT (child_id, source, source_id) DO NOTHING
`,
	tx.ID,
	tx.ChildID,
	tx.Source,
	tx.SourceID,
	tx.XPAmount,
)

	return err
}
func (r *Repo) UpsertXP(
	ctx context.Context,
	childID string,
	amount int,
) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO child_xp (child_id, total_xp, level, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (child_id)
		DO UPDATE SET
			total_xp = child_xp.total_xp + EXCLUDED.total_xp,
			level = ((child_xp.total_xp + EXCLUDED.total_xp) / 100) + 1,
			updated_at = NOW()
	`,
		childID,
		amount,
		(amount/100)+1,
	)

	return err
}
func (r *Repo) GetState(ctx context.Context, childID string) (rules.State, error) {
	var state rules.State

	xp, err := r.GetXP(ctx, childID)
	if err != nil {
		// 👇 THIS IS THE FIX
		if errors.Is(err, sql.ErrNoRows) {
			return rules.State{
				TotalXP: 0,
				Level:   1,
			}, nil
		}
		return state, err
	}

	state.TotalXP = xp.TotalXP
	state.Level = xp.Level

	return state, nil
}
func (r *Repo) GetCurrentPage(
	ctx context.Context,
	childID string,
	bookID string,
) (int, error) {

	var page int

	err := r.db.QueryRow(ctx, `
		SELECT current_page
		FROM book_progress
		WHERE child_id = $1
		  AND book_id = $2
	`,
		childID,
		bookID,
	).Scan(&page)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	return page, err
}
func (r *Repo) GetBadges(
	ctx context.Context,
	childID string,
) ([]model.MilestoneAward, error) {

	rows, err := r.db.Query(ctx, `
		SELECT id, child_id, book_id, milestone_id, reward_type, reward_ref, event_id, created_at
		FROM milestones_awarded
		WHERE child_id = $1
		ORDER BY created_at DESC
	`, childID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.MilestoneAward

	for rows.Next() {
		var m model.MilestoneAward

		err := rows.Scan(
			&m.ID,
			&m.ChildID,
			&m.BookID,
			&m.MilestoneID,
			&m.RewardType,
			&m.RewardRef,
			&m.EventID,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}
