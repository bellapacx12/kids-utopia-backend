package repository

import (
	"context"

	"github.com/bellapacx/kids-utopia/internal/users/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}
func (r *userRepo) FindByID(ctx context.Context, id string) (*model.User, error) {

	query := `
		SELECT id, email, phone, password_hash,
		       name, avatar_url,
		       role, is_verified, is_active,
		       created_at, updated_at
		FROM users
		WHERE id=$1
	`

	var u model.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.PasswordHash,
		&u.Name,
		&u.AvatarURL,
		&u.Role,
		&u.IsVerified,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
func (r *userRepo) Update(ctx context.Context, u *model.User) error {

	query := `
		UPDATE users
		SET name=$1,
		    avatar_url=$2,
		    email=$3,
		    phone=$4,
		    updated_at=NOW()
		WHERE id=$5
	`

	_, err := r.db.Exec(ctx, query,
		u.Name,
		u.AvatarURL,
		u.Email,
		u.Phone,
		u.ID,
	)

	return err
}