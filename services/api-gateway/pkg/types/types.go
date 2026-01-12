package types

type CreateProductRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description" binding:"required"`
	Price       float64           `json:"price" binding:"required,gt=0"`
	CategoryID  string            `json:"category_id" binding:"required"`
	Attributes  map[string]string `json:"attributes" binding:"min=1,max=10"`
	Stock       int               `json:"stock" binding:"required,gte=0"`
}
type CreateProductResponse struct {
	ID string `json:"id"`
}

// t-shirt-basic
// TSB-BLK-M

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

type RotateTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	UserId       string `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
}

// func (req *CreateProductRequest) Validate() error {
//     // 1. ตรวจสอบว่าต้องมี Key ที่จำเป็น (เช่น color, size)
//     requiredKeys := []string{"color", "size"}
//     for _, key := range requiredKeys {
//         if _, ok := req.Attributes[key]; !ok {
//             return fmt.Errorf("attribute '%s' is required", key)
//         }
//     }

//     // 2. ตรวจสอบว่าค่า (Value) ใน Map ห้ามเป็นค่าว่าง
//     for key, value := range req.Attributes {
//         if strings.TrimSpace(value) == "" {
//             return fmt.Errorf("value for attribute '%s' cannot be empty", key)
//         }
//     }

//     return nil
// }
