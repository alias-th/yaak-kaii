package service

import (
	"yaak-kaii/services/auth-service/internal/auth"
	"yaak-kaii/services/auth-service/internal/domain"
)

type service struct {
	userRepo         domain.UserRepository
	guestRepo        domain.GuestRepository
	refreshTokenRepo domain.RefreshTokenRepository
	roleRepo         domain.RoleRepository
	shopRepo         domain.ShopRepository
	jwtAuth          *auth.JWTAuthenticator
}

func NewService(
	userRepo domain.UserRepository,
	guestRepo domain.GuestRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	roleRepo domain.RoleRepository,
	shopRepo domain.ShopRepository,
	jwtAuth *auth.JWTAuthenticator,

) *service {
	return &service{
		userRepo:         userRepo,
		guestRepo:        guestRepo,
		refreshTokenRepo: refreshTokenRepo,
		roleRepo:         roleRepo,
		shopRepo:         shopRepo,
		jwtAuth:          jwtAuth,
	}
}
