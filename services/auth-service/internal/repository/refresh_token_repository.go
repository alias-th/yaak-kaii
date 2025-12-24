package repository

import (
	"context"
	"log"
	"time"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type refreshTokenRepository struct {
	store db.Store
}

func NewRefreshTokenRepository(store db.Store) domain.RefreshTokenRepository {
	return &refreshTokenRepository{store: store}
}

func (r *refreshTokenRepository) CreateRefreshToken(
	ctx context.Context,
	token *domain.RefreshTokenModel) error {
	arg := db.CreateRefreshTokenParams{
		UserID:    pgtype.UUID{Bytes: token.User.ID, Valid: true},
		TokenHash: token.TokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Unix(token.ExpiresAt, 0), Valid: true},
	}
	_, err := r.store.CreateRefreshToken(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}

func (r *refreshTokenRepository) GetRefreshTokenByHash(ctx context.Context, hash string) (*domain.RefreshTokenModel, error) {
	result, err := r.store.GetRefreshTokenByTokenHash(ctx, hash)
	if err != nil {
		log.Printf("failed to get refresh token: error=%v", err)
		if err == pgx.ErrNoRows {
			return nil, utils.NewTokenNotFoundError()
		}

		return nil, err
	}
	return mapDBRefreshTokenToDomain(result), nil
}

func (r *refreshTokenRepository) RevokedRefreshToken(ctx context.Context, id uuid.UUID) error {
	err := r.store.RevokeRefreshToken(ctx, id)
	if err != nil {
		log.Printf("failed to revokes refresh token: error=%v", err)
		return err
	}
	return nil
}

func mapDBRefreshTokenToDomain(dbToken db.RefreshToken) *domain.RefreshTokenModel {
	var revokedAt int64 = 0
	if dbToken.RevokedAt.Valid {
		revokedAt = dbToken.RevokedAt.Time.Unix()
	}

	return &domain.RefreshTokenModel{
		ID: dbToken.ID,
		User: domain.UserModel{
			ID: dbToken.UserID.Bytes,
		},
		CreatedAt: dbToken.CreatedAt.Unix(),
		ExpiresAt: dbToken.ExpiresAt.Time.Unix(),
		RevokedAt: revokedAt,
	}
}
