package domain

type LoginRequest struct {
	Email    string
	Password string
}

type CreateUserRequest struct {
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	Password    string
}

type CreateGuestRequest struct {
	IpAddress string
	UserAgent string
}
