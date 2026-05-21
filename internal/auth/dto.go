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
type ForgotPasswordRequest struct {
	Identifier string `json:"identifier"`
}

type VerifyResetOTPRequest struct {
	Identifier string `json:"identifier"`
	Code       string `json:"code"`
}

type ResetPasswordRequest struct {
	Identifier string `json:"identifier"`
	NewPassword string `json:"new_password"`
	Code        string `json:"code"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
type ResendOTPRequest struct {
	Identifier string `json:"identifier"`
}
type VerificationSessionResponse struct {
	Verified       bool `json:"verified"`
	EmailVerified  bool `json:"email_verified"`
	PhoneVerified  bool `json:"phone_verified"`
}