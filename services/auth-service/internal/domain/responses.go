package domain

type LoginResponse struct {
	UserID       string
	Token        string
	RefreshToken string
	ExpiresAt    int64
}

type CreateUserResponse struct {
	User *UserModel
}

type RotateRefreshTokenResponse struct {
	UserID       string
	Token        string
	RefreshToken string
	ExpiresAt    int64
}
