package domain

import (
	"time"

	"github.com/google/uuid"
)

type GuestModel struct {
	ID        uuid.UUID     `json:"id"`
	TokenHash string        `json:"token_hash"`
	IpAddress string        `json:"ip_address"`
	UserAgent string        `json:"user_agent"`
	CreatedAt int64         `json:"created_at"`
	ExpiresAt int64         `json:"expires_at"`
	MetaData  MetaDataGuest `json:"meta_data"`
}

type MetaDataGuest struct {
	LastActivityTime time.Time `json:"last_activity_time"`
}
