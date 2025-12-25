package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
)

func (s *service) VerifyRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenModel, error) {
	// 1. Get refresh token
	raw, err := utils.DecodeBase64URL(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}
	hashed := utils.HashTokenBytes(raw)
	refreshToken, err := s.refreshTokenRepo.GetRefreshTokenByHash(ctx, hashed)

	// 2. Validate refresh token
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

	// 3. Validate user
	user, err := s.userRepo.GetUserByID(ctx, refreshToken.User.ID)
	if err != nil {
		if utils.IsUserNotFoundError(err) {
			return nil, utils.NewUserNotFoundError()
		}
		return nil, utils.NewInternalServerError()
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
	refreshToken, err := s.VerifyRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	//   2. Generate new JWT access token (1 hour)
	newToken, err := s.jwtAuth.GenerateToken(refreshToken.User.ID)
	if err != nil {
		log.Printf("failed to generate token: error=%v", err)
		return nil, utils.NewInternalServerError()
	}

	//   3. Generate new refresh token (30 days)
	newRefreshToken, newRefreshTokenHashed, err := utils.GenerateTokenPair(32)
	if err != nil {
		log.Printf("failed to generate refresh token: error=%v", err)
		return nil, utils.NewInternalServerError()
	}

	//   4. Save new token to DB
	//   5. Revoke old token
	expiresAt := time.Now().Add(time.Hour * 24 * 30).Unix()
	arg := &domain.RefreshTokenModel{
		User: domain.UserModel{
			ID: refreshToken.User.ID,
		},
		TokenHash: newRefreshTokenHashed,
		ExpiresAt: expiresAt,
	}
	err = s.refreshTokenRepo.RotateRefreshTokenTx(ctx, refreshToken.ID, arg)
	if err != nil {
		log.Printf("failed to rotate refresh token: error=%v", err)
		return nil, utils.NewInternalServerError()
	}

	//   6. Return both new tokens
	return &domain.RotateRefreshTokenResponse{
		UserID:       refreshToken.User.ID.String(),
		Token:        newToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
