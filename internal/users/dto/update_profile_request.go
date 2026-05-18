package dto

type UpdateProfileRequest struct {
	Name      *string `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}