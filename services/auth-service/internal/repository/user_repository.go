package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"

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
		if utils.IsUserAlreadyExistsError(err) {
			log.Printf("duplicate email: %s", user.Email)
			return nil, err
		}

		if strings.Contains(err.Error(), "fk_role") {
			log.Printf("role not found: %s", user.Role.ID)
			return nil, err
		}

		log.Printf("failed to create user: error=%v", err)
		return nil, err
	}
	return mapDBUserToDomain(result), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserModel, error) {
	result, err := r.store.GetUserById(ctx, id)
	if err != nil {
		log.Printf("failed to get user: id=%s, error=%v", id, err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapDBUserToDomain(result), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserModel, error) {
	result, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("failed to get user: email=%s, error=%v", email, err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return mapDBUserToDomain(result), nil
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
