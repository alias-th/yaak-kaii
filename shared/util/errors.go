package util

import "errors"

var (
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenExpired  = errors.New("refresh token expired")
	ErrTokenRevoked  = errors.New("refresh token revoked")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserInactive  = errors.New("user is inactive")
)
