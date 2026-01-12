package domain

import (
	"github.com/google/uuid"
)

type UserModel struct {
	ID            uuid.UUID
	Email         string
	EmailVerified bool
	PasswordHash  string
	FirstName     string
	LastName      string
	PhoneNumber   string
	IsActive      bool
	CreatedAt     int64
	UpdatedAt     int64
	DeletedAt     int64
	Role          RoleModel
}
