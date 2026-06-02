package model

import "time"

type ReadingStreak struct {
	ChildID       string    `db:"child_id" json:"child_id"`
	CurrentStreak int       `db:"current_streak" json:"current_streak"`
	LongestStreak int       `db:"longest_streak" json:"longest_streak"`
	LastReadDate  time.Time `db:"last_read_date" json:"last_read_date"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}