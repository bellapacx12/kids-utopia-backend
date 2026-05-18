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

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}

	var email string
	var phone string

	identifier := strings.TrimSpace(req.Identifier)

	// Detect email vs phone
	if strings.Contains(identifier, "@") {
		email = strings.ToLower(identifier)
	} else {
		phone = identifier
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

	// Generate OTP
	code := generateOTP()

	// Store OTP in Redis
	StoreOTP(identifier, code)

	// TEMP ONLY
	fmt.Println("OTP:", code)
	// REAL EMAIL SEND (SES)
	err = s.otpService.Send(identifier, code)
if err != nil {
	return err
}

	return nil
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