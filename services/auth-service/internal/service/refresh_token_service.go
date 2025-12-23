package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
)

func (s *service) VerifyRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenModel, error) {
	raw, err := utils.DecodeBase64URL(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}
	hashed := utils.HashTokenBytes(raw)
	refreshToken, err := s.refreshTokenRepo.GetRefreshTokenByHash(ctx, hashed)
	if err != nil {
		if utils.IsTokenNotFoundError(err) {
			return nil, err
		}

		return nil, utils.NewInternalServerError()
	}
	if refreshToken.RevokedAt != 0 {
		return nil, utils.NewTokenRevokedError()
	}
	if time.Now().After(time.Unix(refreshToken.ExpiresAt, 0)) {
		return nil, utils.NewTokenExpiredError()
	}

	user, err := s.userRepo.GetUserByID(ctx, refreshToken.User.ID)
	if err != nil {
		return nil, utils.NewInternalServerError()
	}
	if user == nil {
		return nil, utils.NewUserNotFoundError()
	}
	if !user.IsActive {
		return nil, utils.NewUserInactiveError()
	}

	return &domain.RefreshTokenModel{
		ID:        refreshToken.ID,
		User:      *user,
		CreatedAt: refreshToken.CreatedAt,
		ExpiresAt: refreshToken.ExpiresAt,
		RevokedAt: refreshToken.RevokedAt,
	}, nil
}

func (s *service) RotateRefreshToken(
	ctx context.Context,
	token string,
) (*domain.RotateRefreshTokenResponse, error) {
	//   1. Verify old token is valid
	_, err := s.VerifyRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	//   2. Generate new JWT access token (1 hour)
	//   3. Generate new refresh token (30 days)
	//   4. Save new token to DB
	//   5. Revoke old token
	//   6. Return both new tokens

	return nil, errors.New("not implemented")
}
