package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bellapacx/kids-utopia/internal/notifications/otp"
	"github.com/bellapacx/kids-utopia/pkg/security"
)

type Service struct {
	repo      *Repository
	otpService *otp.Service
	jwtSecret string
}

func NewService(repo *Repository, otpService *otp.Service,secret string,) *Service {
	return &Service{
		repo:      repo,
		otpService: otpService,
		jwtSecret: secret,
	}
}

// ========================================
// REGISTER
// ========================================

func (s *Service) Register(req RegisterRequest) error {

	if req.Name == "" {
		return errors.New("name is required")
	}

	email, phone, err := normalizeIdentifier(req.Identifier)
	if err != nil {
		return err
	}

	// 🔥 CHECK IF USER EXISTS FIRST
	exists, err := s.repo.ExistsByIdentifier(context.Background(), email, phone)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(
		context.Background(),
		req.Name,
		email,
		phone,
		hash,
	)
	if err != nil {
		return err
	}

	code := generateOTP()

	StoreOTP(req.Identifier, code)

	fmt.Println("OTP:", code)

	return s.otpService.Send(req.Identifier, code)
}
// ========================================
// LOGIN
// ========================================

func (s *Service) Login(req LoginRequest, deviceID string) (*LoginResponse, error) {

	identifier := strings.TrimSpace(req.Identifier)

	user, err := s.repo.FindByIdentifier(
		context.Background(),
		identifier,
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// BLOCK UNVERIFIED USERS
	if !user.IsVerified {
		return nil, errors.New("account not verified")
	}

	// PASSWORD CHECK
	if !security.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// ACCESS TOKEN
	accessToken, err := security.GenerateToken(
		user.ID,
		user.Role,
		s.jwtSecret,
	)
	if err != nil {
		return nil, err
	}

	// REFRESH TOKEN (RAW - RETURNED TO CLIENT ONLY)
	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// HASH BEFORE DB STORAGE (IMPORTANT FIX)
	tokenHash := security.HashToken(refreshToken)

	// DEFAULT DEVICE
	if deviceID == "" {
		deviceID = "unknown-device"
	}

	// STORE REFRESH TOKEN IN DB
	err = s.repo.StoreRefreshToken(
		context.Background(),
		user.ID,
		tokenHash,
		deviceID,
	)
	if err != nil {
		return nil, err
	}

	// STORE SESSION IN REDIS (FIXED)
	// STORE SESSION (Redis) 
	err = StoreRefreshSession(user.ID, refreshToken) 
	if err != nil { return nil, err }

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
// ========================================
// VERIFY OTP
// ========================================

func (s *Service) VerifyOTP(req VerifyOTPRequest) error {

	ok := s.otpService.Verify(req.Identifier, req.Code)

	if !ok {
		return errors.New("invalid otp")
	}

	return s.repo.VerifyUser(
		context.Background(),
		req.Identifier,
	)
}
func (s *Service) RefreshToken(oldToken string) (*LoginResponse, error) {

	userID, err := s.repo.ValidateRefreshToken(context.Background(), oldToken)
	if err != nil {
		return nil, err
	}

	// REVOKE OLD TOKEN
	s.repo.RevokeToken(context.Background(), oldToken)

	// GENERATE NEW TOKENS
	newAccessToken, _ := security.GenerateToken(userID, "parent", s.jwtSecret)
	newRefreshToken, err := security.GenerateRandomToken()
if err != nil {
	return nil, err
}

	// STORE NEW TOKEN
	s.repo.StoreRefreshToken(context.Background(), userID, newRefreshToken, "device")

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
func normalizeIdentifier(identifier string) (email, phone string, err error) {
	identifier = strings.TrimSpace(identifier)

	if identifier == "" {
		return "", "", errors.New("identifier required")
	}

	if strings.Contains(identifier, "@") {
		email = strings.ToLower(identifier)
		return email, "", nil
	}

	phone = identifier
	return "", phone, nil
}
func (s *Service) ForgotPassword(req ForgotPasswordRequest) error {

	user, err := s.repo.FindByIdentifier(context.Background(), req.Identifier)
	if err != nil || user == nil {
		return errors.New("invalid request")
	}

	code := generateOTP()

	StoreOTP(req.Identifier, code)

	StoreResetSession(req.Identifier)

	return s.otpService.Send(req.Identifier, code)
}
func (s *Service) VerifyResetOTP(req VerifyResetOTPRequest) error {

	ok := s.otpService.Verify(req.Identifier, req.Code)
	if !ok {
		return errors.New("invalid otp")
	}

	return nil
}
func (s *Service) ResetPassword(req ResetPasswordRequest) error {

	valid, err := ValidateResetSession(req.Identifier)
if err != nil || !valid {
	return errors.New("reset session expired")
}

	ok := s.otpService.Verify(req.Identifier, req.Code)
	if !ok {
		return errors.New("invalid otp")
	}

	hash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(context.Background(), req.Identifier, hash)
	if err != nil {
		return err
	}

	// 🔥 invalidate session after success
	DeleteResetSession(req.Identifier)

	return nil
}