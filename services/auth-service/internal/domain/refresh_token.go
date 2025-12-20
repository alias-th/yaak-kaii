package domain

import (
	db "yaak-kaii/services/auth-service/internal/db/sqlc"

	"github.com/google/uuid"
)

type RefreshTokenModel struct {
	ID        uuid.UUID `json:"id"`
	User      db.User   `json:"user"`
	TokenHash string    `json:"token_hash"`
	CreatedAt int64     `json:"created_at"`
	ExpiresAt int64     `json:"expires_at"`
	RevokedAt int64     `json:"revoked_at"`
}
