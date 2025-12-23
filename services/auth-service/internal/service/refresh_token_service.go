package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
)

func (s *service) VerifyRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenModel, error) {
	// decode token
	raw, err := utils.DecodeBase64URL(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}
	// hash token
	hashed := utils.HashTokenBytes(raw)

	// get token
	refreshToken, err := s.refreshTokenRepo.GetRefreshTokenByHash(ctx, hashed)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// check if token expired
	if time.Now().After(time.Unix(refreshToken.ExpiresAt, 0)) {
		return nil, utils.NewTokenExpiredError()
	}

	// check if token revoked
	if refreshToken.RevokedAt != 0 {
		return nil, utils.NewTokenRevokedError()
	}

	// verify user
	user, err := s.userRepo.GetUserByID(ctx, refreshToken.User.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewUserNotFoundError()
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, utils.NewUserInactiveError()
	}

	return &domain.RefreshTokenModel{
		ID:        refreshToken.ID,
		User:      db.User{},
		CreatedAt: refreshToken.CreatedAt,
		ExpiresAt: refreshToken.ExpiresAt,
		RevokedAt: refreshToken.RevokedAt,
	}, nil
}
