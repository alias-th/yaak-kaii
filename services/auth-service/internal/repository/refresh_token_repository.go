package repository

import (
	"context"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
)

type refreshTokenRepository struct {
	store db.Store
}

func NewRefreshTokenRepository(store db.Store) domain.RefreshTokenRepository {
	return &refreshTokenRepository{store: store}
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshTokenModel) error {
	// TODO: Implement
	return nil
}

func (r *refreshTokenRepository) GetRefreshTokenByHash(ctx context.Context, hash string) (*domain.RefreshTokenModel, error) {
	result, err := r.store.GetRefreshTokenByTokenHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return mapDBRefreshTokenToDomain(result), nil
}

func mapDBRefreshTokenToDomain(dbToken db.RefreshToken) *domain.RefreshTokenModel {
	return &domain.RefreshTokenModel{
		ID:        dbToken.ID,
		CreatedAt: dbToken.CreatedAt.Unix(),
		ExpiresAt: dbToken.ExpiresAt.Time.Unix(),
		RevokedAt: dbToken.RevokedAt.Time.Unix(),
	}
}
