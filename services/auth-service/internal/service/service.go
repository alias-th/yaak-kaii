package service

import (
	"yaak-kaii/services/auth-service/internal/domain"
)

type service struct {
	userRepo         domain.UserRepository
	guestRepo        domain.GuestRepository
	refreshTokenRepo domain.RefreshTokenRepository
	roleRepo         domain.RoleRepository
}

func NewService(
	userRepo domain.UserRepository,
	guestRepo domain.GuestRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	roleRepo domain.RoleRepository,
) domain.AuthService {
	return &service{
		userRepo:         userRepo,
		guestRepo:        guestRepo,
		refreshTokenRepo: refreshTokenRepo,
		roleRepo:         roleRepo,
	}
}
