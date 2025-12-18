package service

import (
	"context"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
)

type service struct {
	store domain.UserRepository
}

func NewService(store db.Store) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	// get role user
	role, err := s.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, err
	}

	arg := db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: string(user.Password.PasswordHash),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PhoneNumber:  user.PhoneNumber,
		RoleID:       role.ID,
	}

	result, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &domain.UserModel{
		ID:            result.ID,
		Email:         result.Email,
		EmailVerified: result.EmailVerified,
		FirstName:     result.FirstName,
		LastName:      result.LastName,
		PhoneNumber:   result.PhoneNumber,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt.Unix(),
		UpdatedAt:     result.UpdatedAt.Unix(),
		Role: domain.RoleModel{
			ID: result.RoleID,
		},
	}, nil
}

func (s *service) GetRoleByName(ctx context.Context, name string) (*domain.RoleModel, error) {
	result, err := s.store.GetRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &domain.RoleModel{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}
