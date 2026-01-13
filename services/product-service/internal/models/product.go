package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShopID      uuid.UUID        `gorm:"column:shop_id;type:uuid;not null;uniqueIndex:ux_products_shop_slug"`
	CategoryID  uuid.UUID        `gorm:"column:category_id;type:uuid;not null;index"`
	Category    Category         `gorm:"foreignKey:CategoryID;references:ID"`
	Slug        string           `gorm:"column:slug;type:varchar(255);not null;uniqueIndex:ux_products_shop_slug"`
	Name        string           `gorm:"column:name;type:varchar(255);not null"`
	Description string           `gorm:"column:description;type:text"`
	Status      string           `gorm:"column:status;type:varchar(50);default:'active';not null"`
	CreatedAt   time.Time        `gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime"`
	Images      []ProductImage   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Variants    []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID uuid.UUID `gorm:"column:product_id;type:uuid;not null"`
	ImageURL  string    `gorm:"column:image_url;type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type ProductVariant struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID  uuid.UUID  `gorm:"column:product_id;type:uuid;not null;uniqueIndex:ux_product_variants_sku"`
	Sku        string     `gorm:"column:sku;type:varchar(100);not null;uniqueIndex:ux_product_variants_sku"`
	Price      float64    `gorm:"type:numeric(12,2);not null"`
	Stock      int        `gorm:"column:stock;type:int;not null"`
	Attributes Attributes `gorm:"column:attributes;type:jsonb;default:'{}'"`
	VariantNo  int        `gorm:"column:variant_no;not null;"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
}

type Attributes []byte

func (a Attributes) Value() (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal(a, &result)
	return result, err
}
