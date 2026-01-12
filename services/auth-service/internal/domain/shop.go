package domain

type ShopModel struct {
	ID          string `json:"id"`
	SellerID    string `json:"seller_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
