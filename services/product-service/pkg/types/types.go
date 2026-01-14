package types

import (
	"encoding/json"
	"errors"
	"yaak-kaii/services/product-service/internal/models"
)

type ListProductsReq struct {
	Q           string
	CategoryIds []string
	MinPrice    int64
	MaxPrice    int64
	Limit       int32
	Page        int32
	Sort        string
}
type ListProductsRes struct {
	Products   []models.Product
	Pagination Pagination
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	LastPage int   `json:"last_page"`
	NextPage *int  `json:"next_page,omitempty"`
	PrevPage *int  `json:"prev_page,omitempty"`
	HasNext  bool  `json:"has_next"`
	HasPrev  bool  `json:"has_prev"`
}

type Attributes map[string]string

func (a Attributes) Value() ([]byte, error) {
	return json.Marshal(a)
}
func (a *Attributes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type CreateProductPayload struct {
	UserId       string  // ID of the user creating the product
	Name         string  // Product name
	Description  string  // Product description
	Price        float64 // Product price
	CategoryId   string  // Product category
	Stock        int32   // Available stock quantity
	Attributes   map[string]string
	CategoryName string
}

type CreateProductVariantPayload struct {
	ProductID  string
	Price      float64
	Stock      int32
	Attributes Attributes
}

type AxisDef struct {
	Key      string
	Required bool
	Options  map[string]struct{} // set
	Order    int
}
