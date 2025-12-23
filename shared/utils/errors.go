package utils

import "strings"

type ErrorCode string

const (
	ErrCodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserExists         ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenRevoked       ErrorCode = "TOKEN_REVOKED"
	ErrCodeUserInactive       ErrorCode = "USER_INACTIVE"
	ErrCodeInvalidToken       ErrorCode = "INVALID_TOKEN"
	ErrCodeRoleNotFound       ErrorCode = "ROLE_NOT_FOUND"
	ErrCodeInvalidPassword    ErrorCode = "INVALID_PASSWORD"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
)

type CustomError struct {
	Code       ErrorCode
	Message    string
	HTTPStatus int
}

func (e *CustomError) Error() string {
	return e.Message
}

// Helper functions
func NewUserNotFoundError() *CustomError {
	return &CustomError{
		Code:       ErrCodeUserNotFound,
		Message:    "user not found",
		HTTPStatus: 404,
	}
}

func NewUserAlreadyExistsError() *CustomError {
	return &CustomError{
		Code:       ErrCodeUserExists,
		Message:    "user already exists",
		HTTPStatus: 409,
	}
}

func NewTokenExpiredError() *CustomError {
	return &CustomError{
		Code:       ErrCodeTokenExpired,
		Message:    "refresh token expired",
		HTTPStatus: 401,
	}
}

func NewRoleNotFoundError() *CustomError {
	return &CustomError{
		Code:       ErrCodeRoleNotFound,
		Message:    "user role not found",
		HTTPStatus: 500,
	}
}

func NewTokenRevokedError() *CustomError {
	return &CustomError{
		Code:       ErrCodeTokenRevoked,
		Message:    "refresh token revoked",
		HTTPStatus: 401,
	}
}

func NewUserInactiveError() *CustomError {
	return &CustomError{
		Code:       ErrCodeUserInactive,
		Message:    "user is inactive",
		HTTPStatus: 403,
	}
}

func NewInvalidPasswordError() *CustomError {
	return &CustomError{
		Code:       ErrCodeInvalidPassword,
		Message:    "password is incorrect",
		HTTPStatus: 401,
	}
}

func NewInvalidCredentialsError() *CustomError {
	return &CustomError{
		Code:       ErrCodeInvalidCredentials,
		Message:    "invalid email or password",
		HTTPStatus: 401,
	}
}

// sql
func IsUserAlreadyExistsError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "email")
}
