package auth

type RegisterRequest struct {
	Identifier string `json:"identifier" validate:"required,min=3"`
	Password   string `json:"password" validate:"required,min=8"`
	Name       string `json:"name" validate:"required,min=4"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
	DeviceID   string `json:"device_id"`
}

type VerifyOTPRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Code       string `json:"code" validate:"required,len=6"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}