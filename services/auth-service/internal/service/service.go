package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"time"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/util"
)

type service struct {
	store domain.AuthRepository
}

func NewService(store db.Store) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateGuest(ctx context.Context, guest *domain.GuestModel) (*domain.GuestModel, error) {

	// add expires date
	now := time.Now().UTC()
	expires30Day := now.AddDate(0, 0, 30)
	guest.MetaData.LastActivityTime = expires30Day

	byteData, err := json.Marshal(guest.MetaData)
	if err != nil {
		return nil, err
	}

	// parse IP address string into netip.Addr
	ip, err := netip.ParseAddr(guest.IpAddress)
	if err != nil {
		return nil, err
	}

	// create guest
	arg := db.CreateGuestParams{
		TokenHash: guest.TokenHash,
		IpAddr:    ip,
		UserAgent: guest.UserAgent,
		Metadata:  byteData,
	}
	res, err := s.store.CreateGuest(ctx, arg)
	if err != nil {
		return nil, err
	}

	// response
	return &domain.GuestModel{
		ID:        res.ID,
		TokenHash: res.TokenHash,
		IpAddress: res.IpAddr.String(),
		UserAgent: res.UserAgent,
		CreatedAt: res.CreatedAt.Unix(),
		ExpiresAt: res.CreatedAt.Unix(),
		MetaData: domain.MetaDataGuest{
			LastActivityTime: guest.MetaData.LastActivityTime,
		},
	}, nil
}

func (s *service) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	// get role user
	role, err := s.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, err
	}

	arg := db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
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

func (s *service) VerifyGuestToken(ctx context.Context, token string) (*domain.GuestModel, error) {
	// decode token
	raw, err := util.DecodeBase64URL(token)
	if err != nil {
		return nil, err
	}

	// hash token
	hashed := util.HashTokenBytes(raw)

	guest, err := s.store.GetGuestByToken(ctx, hashed)
	if err != nil {
		return nil, err
	}

	// verify token
	ok, err := util.VerifyToken(token, guest.TokenHash)
	if err != nil || !ok {
		return nil, err
	}

	return &domain.GuestModel{
		ID:        guest.ID,
		IpAddress: guest.IpAddr.String(),
		UserAgent: guest.UserAgent,
		CreatedAt: guest.CreatedAt.Unix(),
		ExpiresAt: guest.CreatedAt.Unix(),
	}, nil

}

func (s *service) VerifyRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenModel, error) {
	// decode token
	raw, err := util.DecodeBase64URL(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}
	// hash token
	hashed := util.HashTokenBytes(raw)

	// get token
	refreshToken, err := s.store.GetRefreshTokenByTokenHash(ctx, hashed)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// check if token expired
	if time.Now().After(refreshToken.ExpiresAt.Time) {
		return nil, util.ErrTokenExpired
	}

	// check if token revoked
	if refreshToken.RevokedAt.Valid {
		return nil, util.ErrTokenRevoked
	}

	// verify user
	user, err := s.store.GetUserById(ctx, refreshToken.UserID.Bytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, util.ErrUserInactive
	}

	return &domain.RefreshTokenModel{
		ID:        refreshToken.ID,
		User:      user,
		CreatedAt: refreshToken.CreatedAt.Unix(),
		ExpiresAt: refreshToken.ExpiresAt.Time.Unix(),
		RevokedAt: refreshToken.RevokedAt.Time.Unix(),
	}, nil

}
