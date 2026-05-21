package model

import "time"

type Subscription struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Plan      string     `db:"plan"`
	Status    string     `db:"status"` // active, expired, canceled

	StartDate time.Time  `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}