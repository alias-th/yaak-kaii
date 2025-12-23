package domain

import (
	"context"
)

type AuthService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	CreateGuest(ctx context.Context, guest *CreateGuestRequest) (string, error)
	VerifyGuestToken(ctx context.Context, token string) (*GuestModel, error)
	VerifyRefreshToken(ctx context.Context, token string) (*RefreshTokenModel, error)
	RotateRefreshToken(ctx context.Context, token string) (*RotateRefreshTokenResponse, error)
}
