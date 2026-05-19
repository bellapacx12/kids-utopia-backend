package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bellapacx/kids-utopia/internal/notifications/otp"
	"github.com/bellapacx/kids-utopia/pkg/redis"
	"github.com/bellapacx/kids-utopia/pkg/security"

	goredis "github.com/redis/go-redis/v9"
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

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {

	if req.Name == "" {
		return errors.New("name is required")
	}

	email, phone, err := normalizeIdentifier(req.Identifier)
	if err != nil {
		return err
	}

	// 🔥 CHECK IF USER EXISTS FIRST
	exists, err := s.repo.ExistsByIdentifier(ctx, email, phone)
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
		ctx,
		req.Name,
		email,
		phone,
		hash,
	)
	if err != nil {
		return err
	}

	code := generateOTP()

	key := strings.ToLower(strings.TrimSpace(req.Identifier))

StoreOTP(key, code)

return s.otpService.Send(key, code)
}
// ========================================
// LOGIN
// ========================================

func (s *Service) Login(ctx context.Context, req LoginRequest, deviceID string) (*LoginResponse, error) {
    
	key := "login:attempt:" + strings.ToLower(strings.TrimSpace(req.Identifier))
	val, err := redis.Client.Get(ctx, key).Int()
if err != nil && err != goredis.Nil {
	return nil, err
}
if val >= 5 {
	return nil, errors.New("too many login attempts, try again later")
}
	identifier := strings.ToLower(strings.TrimSpace(req.Identifier))

    user, err := s.repo.FindByIdentifier(ctx, identifier)

	
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
	redis.Client.Incr(ctx, key)
redis.Client.Expire(ctx, key, 10*time.Minute)

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
		ctx,
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
    redis.Client.Del(ctx, key)
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

	
}
// ========================================
// VERIFY OTP
// ========================================

func (s *Service) VerifyOTP(ctx context.Context, req VerifyOTPRequest) error {

key := strings.ToLower(strings.TrimSpace(req.Identifier))

ok := s.otpService.Verify(key, req.Code)
	if !ok {
		return errors.New("invalid otp")
	}

	return s.repo.VerifyUser(
		ctx,
		req.Identifier,
	)
}
func (s *Service) RefreshToken(ctx context.Context, oldToken string) (*LoginResponse, error) {

	// 1. validate old token
	userID, err := s.repo.ValidateRefreshToken(ctx, oldToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. revoke old token in DB
	err = s.repo.RevokeToken(ctx, oldToken)
	if err != nil {
		return nil, err
	}

	// 3. remove old Redis session
	_ = DeleteRefreshSession(oldToken)

	// 4. generate new tokens
	newAccessToken, err := security.GenerateToken(userID, "parent", s.jwtSecret)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 5. store new refresh token in DB
	err = s.repo.StoreRefreshToken(ctx, userID, newRefreshToken, "device")
	if err != nil {
		return nil, err
	}

	// 6. store new session in Redis
	err = StoreRefreshSession(userID, newRefreshToken)
	if err != nil {
		return nil, err
	}

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
func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {

	user, err := s.repo.FindByIdentifier(ctx, req.Identifier)
	if err != nil || user == nil {
		return errors.New("invalid request")
	}

	code := generateOTP()

	StoreOTP(req.Identifier, code)

	StoreResetSession(req.Identifier)

	return s.otpService.Send(req.Identifier, code)
}
func (s *Service) VerifyResetOTP(ctx context.Context,req VerifyResetOTPRequest) error {

	ok := s.otpService.Verify(req.Identifier, req.Code)
	if !ok {
		return errors.New("invalid otp")
	}

	return nil
}
func (s *Service) ResetPassword(ctx context.Context,req ResetPasswordRequest) error {

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

	err = s.repo.UpdatePassword(ctx, req.Identifier, hash)
	if err != nil {
		return err
	}

	// 🔥 invalidate session after success
	DeleteResetSession(req.Identifier)

	return nil
}
func (s *Service) Logout(ctx context.Context, refreshToken string) error {

	// revoke token in postgres
	err := s.repo.RevokeToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	// remove redis session
	err = DeleteRefreshSession(refreshToken)
	if err != nil {
		return err
	}

	return nil
}