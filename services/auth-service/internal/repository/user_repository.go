package repository

import (
	"context"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"

	"github.com/google/uuid"
)

type userRepository struct {
	store db.Store
}

func NewUserRepository(store db.Store) domain.UserRepository {
	return &userRepository{store: store}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	arg := db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PhoneNumber:  user.PhoneNumber,
		RoleID:       user.Role.ID,
	}
	result, err := r.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return mapDBUserToDomain(result), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserModel, error) {
	result, err := r.store.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapDBUserToDomain(result), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserModel, error) {
	// TODO: Implement if needed
	return nil, nil
}

func mapDBUserToDomain(dbUser db.User) *domain.UserModel {
	return &domain.UserModel{
		ID:            dbUser.ID,
		Email:         dbUser.Email,
		EmailVerified: dbUser.EmailVerified,
		PasswordHash:  dbUser.PasswordHash,
		FirstName:     dbUser.FirstName,
		LastName:      dbUser.LastName,
		PhoneNumber:   dbUser.PhoneNumber,
		IsActive:      dbUser.IsActive,
		CreatedAt:     dbUser.CreatedAt.Unix(),
		UpdatedAt:     dbUser.UpdatedAt.Unix(),
		Role: &domain.RoleModel{
			ID: dbUser.RoleID,
		},
	}
}
