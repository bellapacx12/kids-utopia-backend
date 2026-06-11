package model

import "time"

type Child struct {
	ID        string
	ParentID  string

	Name      string
	AvatarURL *string
	Age       int
	Language  string
  
	CreatedAt time.Time
	UpdatedAt time.Time
}