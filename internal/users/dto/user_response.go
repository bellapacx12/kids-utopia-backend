package dto

type UserResponse struct {
	ID         string  `json:"id"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`

	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`

	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
	IsActive   bool   `json:"is_active"`
}
