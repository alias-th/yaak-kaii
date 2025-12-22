package types

type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	FirstName   string `json:"first_name" binding:"required,min=2,max=50"`
	LastName    string `json:"last_name" binding:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" binding:"required,numeric,len=10"`
	Password    string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"`
}

type CreateGuestResponse struct {
	GuestToken string `json:"guest_token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserId       string `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
}
