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
