package types

import (
	"time"

	"github.com/google/uuid"
)

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

type ListProductRequest struct {
	Q          string   `form:"q"`
	CategoryID []string `form:"category_id"`
	MinPrice   *int64   `form:"min_price"`
	MaxPrice   *int64   `form:"max_price"`
	Page       int32    `form:"page,default=1"`
	Limit      int32    `form:"limit,default=20"`
	Sort       string   `form:"sort,default=newest"`
}

type ListProductResponse struct {
	Products   []ProductResponse `json:"products"`
	Pagination Pagination        `json:"pagination"`
}

type ProductResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	ShopID      string           `json:"shop_id"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Variants    []ProductVariant `json:"variants"`
	Images      []string         `json:"images"`
	Category    ProductCategory  `json:"category"`
}

type ProductVariant struct {
	SKU        string            `json:"sku"`
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Attributes map[string]string `json:"attributes"`
}

type ProductCategory struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
}

type Pagination struct {
	Page     int32 `json:"page"`
	Limit    int32 `json:"limit"`
	Total    int32 `json:"total"`
	LastPage int32 `json:"last_page"`
	NextPage int32 `json:"next_page"`
	PrevPage int32 `json:"prev_page"`
	HasNext  bool  `json:"has_next"`
	HasPrev  bool  `json:"has_prev"`
}

type CreateProductVariantRequest struct {
	Price      float64           `json:"price" binding:"required,gt=0"`
	Stock      int32             `json:"stock" binding:"required,gte=0"`
	Attributes map[string]string `json:"attributes" binding:"min=1,max=10"`
}

type CreateProductVariantResponse struct {
	VariantID string `json:"variant_id"`
}

type ProductDetailResponse struct {
	Product struct {
		ID          uuid.UUID `json:"id"`
		ShopID      uuid.UUID `json:"shop_id"`
		CategoryID  uuid.UUID `json:"category_id"`
		Name        string    `json:"name"`
		Slug        string    `json:"slug"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"product"`

	Axes     []AxisDTO    `json:"axes"`
	Variants []VariantDTO `json:"variants"`
	Images   []string     `json:"images"`
}

type AxisDTO struct {
	Key      string   `json:"key"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Options  []string `json:"options"`
	Order    int      `json:"order"`
}

type VariantDTO struct {
	ID         uuid.UUID         `json:"id"`
	SKU        string            `json:"sku"`
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Attributes map[string]string `json:"attributes"`
	VariantKey string            `json:"variant_key,omitempty"`
}
