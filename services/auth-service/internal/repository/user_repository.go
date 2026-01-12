package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

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

func (r *userRepository) CreateSellerTx(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	// payload
	arg := db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PhoneNumber:  user.PhoneNumber,
	}
	roleId := user.Role.ID

	// result
	var resultUser db.User

	// tx
	err := r.store.WithTx(ctx, func(q *db.Queries) error {

		// create user
		createdUser, err := q.CreateUser(ctx, arg)
		if err != nil {
			log.Printf("failed to create user: error=%v", err)
			return err
		}

		// create seller
		createdSeller, err := q.CreateSeller(ctx, db.CreateSellerParams{UserID: createdUser.ID, Name: createdUser.FirstName})
		if err != nil {
			log.Printf("failed to create seller: error=%v", err)
			return err
		}

		// create shop
		_, err = q.CreateShop(ctx, db.CreateShopParams{SellerID: createdSeller.ID,
			Name:        createdSeller.Name,
			Slug:        slug.Make(createdUser.ID.String() + " " + createdUser.FirstName),
			Description: pgtype.Text{Valid: true, String: ""},
		})
		if err != nil {
			log.Printf("failed to create shop: error=%v", err)
			return err
		}

		// create user&role
		_, err = q.CreateUserRole(ctx, db.CreateUserRoleParams{UserID: createdUser.ID, RoleID: roleId})
		if err != nil {
			log.Printf("failed to create user&role: error=%v", err)
			return err
		}

		resultUser = createdUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return mapDBUserToDomain(resultUser), nil
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
	}
}
