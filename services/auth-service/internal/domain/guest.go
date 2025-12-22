package domain

import (
	"net/netip"

	"github.com/google/uuid"
)

type GuestModel struct {
	ID        uuid.UUID     `json:"id"`
	TokenHash string        `json:"token_hash"`
	IpAddress netip.Addr    `json:"ip_address"`
	UserAgent string        `json:"user_agent"`
	CreatedAt int64         `json:"created_at"`
	ExpiresAt int64         `json:"expires_at"`
	MetaData  MetaDataGuest `json:"meta_data"`
}

type MetaDataGuest struct {
	LastActivityTime int64  `json:"last_activity_time"`
	SessionStartTime int64  `json:"session_start_time"`
	Os               string `json:"os"`
	UserAgent        string `json:"user_agent"`
	Device           string `json:"device"`
	RequestCount     int    `json:"request_count"`
}
