package domain

import (
	"context"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserModel struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	PasswordHash  string    `json:"-"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	PhoneNumber   string    `json:"phone_number"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     int64     `json:"created_at"`
	UpdatedAt     int64     `json:"updated_at"`
	DeletedAt     int64     `json:"deleted_at,omitempty"`
	Role          RoleModel `json:"role"`
}

type AuthRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetRoleByName(ctx context.Context, name string) (db.Role, error)
	CreateGuest(ctx context.Context, arg db.CreateGuestParams) (db.Guest, error)
	GetGuestByToken(ctx context.Context, token string) (db.Guest, error)
	GetRefreshTokenByTokenHash(ctx context.Context, token string) (db.RefreshToken, error)
	GetUserById(ctx context.Context, id uuid.UUID) (db.User, error)
}

type AuthService interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	GetRoleByName(ctx context.Context, name string) (*RoleModel, error)
	CreateGuest(ctx context.Context, guest *GuestModel) (*GuestModel, error)
	VerifyGuestToken(ctx context.Context, token string) (*GuestModel, error)
	VerifyRefreshToken(ctx context.Context, token string) (*RefreshTokenModel, error)
}
