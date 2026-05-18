package model

import "time"

type User struct {
	ID           string
	Email        *string
	Phone        *string
	PasswordHash string

	Name      *string
	AvatarURL *string

	Role       string
	IsVerified bool
	IsActive   bool

	CreatedAt time.Time
	UpdatedAt time.Time
}