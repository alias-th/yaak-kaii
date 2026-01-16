package types

type CreateCategoryRequest struct {
	Name        string              `json:"name" binding:"required,min=2,max=255"`
	Description string              `json:"description" binding:"required,min=2,max=255"`
	Attributes  []CategoryAttribute `json:"attributes" binding:"omitempty,dive"`
}
type CreateCategoryResponse struct {
	CategoryID string `json:"category_id"`
}

type CategoryAttribute struct {
	// CategoryId string   `json:"category_id" binding:"required"`
	Key       string   `json:"key" binding:"required"`
	Label     string   `json:"label" binding:"required"`
	Type      string   `json:"type" binding:"required,min=2,max=50"`
	Required  bool     `json:"required"`
	Options   []string `json:"options" binding:"omitempty,min=1,dive,required,min=1,max=100"`
	Scope     string   `json:"scope" binding:"required,min=2,max=50,oneof=PRODUCT VARIANT"`
	AxisOrder int64    `json:"axis_order" binding:"required"`
}
