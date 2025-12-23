package domain

import (
	"time"
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

// Expired2Day returns the time when token expires (2 days from now)
func (reft *RefreshTokenModel) Expired2Day() time.Time {
	return time.Now().Add(2 * 24 * time.Hour)
}

// IsExpired checks if token has expired
func (reft *RefreshTokenModel) IsExpired() bool {
	return time.Now().Unix() > reft.ExpiresAt
}

// IsRevoked checks if token has been revoked
func (reft *RefreshTokenModel) IsRevoked() bool {
	return reft.RevokedAt > 0
}

// IsValid checks if token is still valid
func (reft *RefreshTokenModel) IsValid() bool {
	return !reft.IsExpired() && !reft.IsRevoked()
}
