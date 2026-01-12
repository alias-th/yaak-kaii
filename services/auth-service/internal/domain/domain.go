package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	CreateSellerTx(ctx context.Context, user *UserModel) (*UserModel, error)
}

type GuestRepository interface {
	CreateGuest(ctx context.Context, guest *GuestModel) (*GuestModel, error)
	GetGuestByToken(ctx context.Context, token string) (*GuestModel, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *RefreshTokenModel) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (*RefreshTokenModel, error)
	RevokedRefreshToken(ctx context.Context, id uuid.UUID) error
	RotateRefreshTokenTx(ctx context.Context, oldTokenID uuid.UUID, newToken *RefreshTokenModel) error
}

type RoleRepository interface {
	GetRoleByName(ctx context.Context, name string) (*RoleModel, error)
}

type ShopRepository interface {
	GetShopByUserID(ctx context.Context, userId uuid.UUID) (*ShopModel, error)
}
