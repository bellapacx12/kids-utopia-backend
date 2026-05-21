package auth

import (
	"context"

	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/security"
)

type Repository struct{}

type User struct {
	ID           string
	Name         string
	Email        *string
	Phone        *string
	PasswordHash string
	Role         string
	IsVerified   bool
	
	EmailVerified   bool
	PhoneVerified   bool
}

func (r *Repository) CreateUser(
	ctx context.Context,
	name,email, phone, passwordHash string,
) error {

	var emailPtr *string
	var phonePtr *string

	if email == "" {
		emailPtr = nil
	} else {
		emailPtr = &email
	}

	if phone == "" {
		phonePtr = nil
	} else {
		phonePtr = &phone
	}

	_, err := database.DB.Exec(ctx,
		`INSERT INTO users (name, email, phone, password_hash)
		 VALUES ($1, $2, $3, $4)`,
		name,
		emailPtr,
		phonePtr,
		passwordHash,
	)

	return err
}
func (r *Repository) FindByIdentifier(ctx context.Context, identifier string) (*User, error) {

	var user User

	err := database.DB.QueryRow(ctx, `
		SELECT id, email, phone, password_hash, role, is_verified, email_verified,
		phone_verified
		FROM users
		WHERE email = $1 OR phone = $1
		LIMIT 1
	`, identifier).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.EmailVerified,
	&user.PhoneVerified,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *Repository) VerifyUser(
	ctx context.Context,
	identifier string,
) error {

	_, err := database.DB.Exec(ctx, `
		UPDATE users
		SET is_verified = true
		WHERE email = $1 OR phone = $1
	`, identifier)

	return err
}
func (r *Repository) StoreRefreshToken(
	ctx context.Context,
	userID string,
	token string,
	deviceID string,
) error {

	tokenHash := security.HashToken(token)

	_, err := database.DB.Exec(ctx, `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			device_id,
			expires_at
		)
		VALUES ($1, $2, $3, NOW() + INTERVAL '7 days')
	`,
		userID,
		tokenHash,
		deviceID,
	)

	return err
}
func (r *Repository) ValidateRefreshToken(
	ctx context.Context,
	token string,
) (string, error) {

	tokenHash := security.HashToken(token)

	var userID string

	err := database.DB.QueryRow(ctx, `
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1
		AND revoked = false
		AND expires_at > NOW()
	`, tokenHash).Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, nil
}
func (r *Repository) RevokeToken(
	ctx context.Context,
	token string,
) error {

	tokenHash := security.HashToken(token)

	_, err := database.DB.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE token_hash = $1
	`, tokenHash)

	return err
}
func (r *Repository) ExistsByIdentifier(
	ctx context.Context,
	email string,
	phone string,
) (bool, error) {

	var exists bool

	err := database.DB.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE
				($1 <> '' AND email = $1)
				OR
				($2 <> '' AND phone = $2)
		)
	`, email, phone).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) UpdatePassword(
	ctx context.Context,
	identifier string,
	passwordHash string,
) error {

	_, err := database.DB.Exec(ctx, `
		UPDATE users
		SET password_hash = $1
		WHERE email = $2 OR phone = $2
	`, passwordHash, identifier)

	return err
}
func (r *Repository) VerifyEmail(
	ctx context.Context,
	email string,
) error {

	var id string

	err := database.DB.QueryRow(ctx, `
		UPDATE users
		SET
			email_verified = TRUE,
			is_verified = TRUE,
			updated_at = NOW()
		WHERE email = $1
		RETURNING id
	`, email).Scan(&id)

	return err
}
func (r *Repository) VerifyPhone(
	ctx context.Context,
	phone string,
) error {

	var id string

	err := database.DB.QueryRow(ctx, `
		UPDATE users
		SET
			phone_verified = TRUE,
			is_verified = TRUE,
			updated_at = NOW()
		WHERE phone = $1
		RETURNING id
	`, phone).Scan(&id)

	return err
}
func (r *Repository) FindByID(
	ctx context.Context,
	userID string,
) (*User, error) {

	var user User

	err := database.DB.QueryRow(ctx, `
		SELECT
			id,
			name,
			email,
			phone,
			password_hash,
			role,
			is_verified,
			email_verified,
			phone_verified
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.EmailVerified,
		&user.PhoneVerified,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}