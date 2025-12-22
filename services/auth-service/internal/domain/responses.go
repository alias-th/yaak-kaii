package domain

type LoginResponse struct {
	User         *UserModel
	Token        string
	RefreshToken string
	ExpiresAt    int64
}

type CreateUserResponse struct {
	User *UserModel
}
