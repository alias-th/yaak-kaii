package domain

import (
	"context"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Password      Password  `json:"-"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	PhoneNumber   string    `json:"phone_number"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     int64     `json:"created_at"`
	UpdatedAt     int64     `json:"updated_at"`
	DeletedAt     int64     `json:"deleted_at,omitempty"`
	Role          RoleModel `json:"role"`
}

type Password struct {
	Password     *string
	PasswordHash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Password = &text
	p.PasswordHash = hash

	return nil
}

func (p *Password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.PasswordHash, []byte(text))
}

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetRoleByName(ctx context.Context, name string) (db.Role, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	GetRoleByName(ctx context.Context, name string) (*RoleModel, error)
}
