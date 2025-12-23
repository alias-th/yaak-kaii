package repository

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
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
		if strings.Contains(err.Error(), "duplicate key") {
			log.Printf("user already exists: email=%s", user.Email)
			return nil, utils.NewUserAlreadyExistsError()
		}
		if strings.Contains(err.Error(), "fk_role") {
			log.Printf("role not found: roleId=%s", user.Role.ID)
			return nil, utils.NewRoleNotFoundError()
		}

		log.Printf("failed to create user: error=%v", err)
		return nil, err
	}
	return mapDBUserToDomain(result), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserModel, error) {
	result, err := r.store.GetUserById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.NewUserNotFoundError()
		}
		log.Printf("failed to get user: error=%v", err)
		return nil, err
	}
	return mapDBUserToDomain(result), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserModel, error) {
	result, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.NewUserNotFoundError()
		}
		log.Printf("failed to get user: error=%v", err)
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
