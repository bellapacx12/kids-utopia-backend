package dto

import "time"

type ChildResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Age       *int    `json:"age,omitempty"`
	Language  string  `json:"language"`

	CreatedAt time.Time `json:"created_at"`
}