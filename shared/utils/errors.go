package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeTokenNotFound      ErrorCode = "TOKEN_NOT_FOUND"
)

type GRPCErrorCode int

const (
	OK                 GRPCErrorCode = 0
	Cancelled          GRPCErrorCode = 1
	Unknown            GRPCErrorCode = 2
	InvalidArg         GRPCErrorCode = 3
	DeadlineExceeded   GRPCErrorCode = 4
	NotFound           GRPCErrorCode = 5
	AlreadyExists      GRPCErrorCode = 6
	PermissionDenied   GRPCErrorCode = 7
	ResourceExhausted  GRPCErrorCode = 8
	FailedPrecondition GRPCErrorCode = 9
	Aborted            GRPCErrorCode = 10
	OutOfRange         GRPCErrorCode = 11
	Unimplemented      GRPCErrorCode = 12
	Internal           GRPCErrorCode = 13
	Unavailable        GRPCErrorCode = 14
	DataLoss           GRPCErrorCode = 15
	Unauthenticated    GRPCErrorCode = 16
)

type CustomError struct {
	Code       ErrorCode
	Message    string
	HTTPStatus int
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewGRPCError creates a gRPC status error from code and message
func NewGRPCError(code GRPCErrorCode, message string) error {
	return status.Error(codes.Code(code), message)
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

func NewInternalServerError() *CustomError {
	return &CustomError{
		Code:       ErrCodeInternalError,
		Message:    "internal server error",
		HTTPStatus: 500,
	}
}

func NewTokenNotFoundError() *CustomError {
	return &CustomError{
		Code:       ErrCodeTokenNotFound,
		Message:    "token not found",
		HTTPStatus: 404,
	}
}

// Error type checkers - Helper functions to check error types
func hasErrorCode(err error, code ErrorCode) bool {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Code == code
	}
	return false
}

// IsUserAlreadyExistsError checks if error is due to user already exists (duplicate email)
func IsUserAlreadyExistsError(err error) bool {
	return hasErrorCode(err, ErrCodeUserExists)
}

// IsTokenNotFoundError checks if error is due to token not found
func IsTokenNotFoundError(err error) bool {
	return hasErrorCode(err, ErrCodeTokenNotFound)
}

// IsTokenExpiredError checks if error is due to token expiration
func IsTokenExpiredError(err error) bool {
	return hasErrorCode(err, ErrCodeTokenExpired)
}

// IsTokenRevokedError checks if error is due to token being revoked
func IsTokenRevokedError(err error) bool {
	return hasErrorCode(err, ErrCodeTokenRevoked)
}

// IsUserNotFoundError checks if error is due to user not found
func IsUserNotFoundError(err error) bool {
	return hasErrorCode(err, ErrCodeUserNotFound)
}

// IsUserInactiveError checks if error is due to user being inactive
func IsUserInactiveError(err error) bool {
	return hasErrorCode(err, ErrCodeUserInactive)
}

// IsRoleNotFoundError checks if error is due to role not found
func IsRoleNotFoundError(err error) bool {
	return hasErrorCode(err, ErrCodeRoleNotFound)
}

// ToGRPCError converts a CustomError to a gRPC status error
// Returns nil if error is not a CustomError
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	customErr, ok := err.(*CustomError)
	if !ok {
		return err
	}

	// Map CustomError codes to gRPC codes
	var grpcCode GRPCErrorCode
	switch customErr.Code {
	case ErrCodeUserNotFound, ErrCodeTokenNotFound, ErrCodeRoleNotFound:
		grpcCode = NotFound
	case ErrCodeUserExists:
		grpcCode = AlreadyExists
	case ErrCodeTokenExpired, ErrCodeTokenRevoked, ErrCodeInvalidToken:
		grpcCode = Unauthenticated
	case ErrCodeInvalidPassword, ErrCodeInvalidCredentials:
		grpcCode = Unauthenticated
	case ErrCodeUserInactive:
		grpcCode = PermissionDenied
	default:
		grpcCode = Internal
	}

	return NewGRPCError(grpcCode, customErr.Message)
}
