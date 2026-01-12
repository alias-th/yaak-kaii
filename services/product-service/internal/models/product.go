package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ShopID      uuid.UUID `gorm:"column:shop_id;type:uuid;not null;index"`
	CategoryID  uuid.UUID `gorm:"column:category_id;type:uuid;not null;index"`
	Slug        string    `gorm:"column:slug;type:varchar(255);not null;uniqueIndex"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Description string    `gorm:"column:description;type:text"`
	Status      string    `gorm:"column:status;type:varchar(50);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Images      []ProductImage   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Variants    []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID uuid.UUID `gorm:"column:product_id;type:uuid;not null"`
	ImageURL  string    `gorm:"column:image_url;type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductVariant struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID uuid.UUID `gorm:"column:product_id;type:uuid;not null"`
	Sku       string    `gorm:"column:sku;type:varchar(100);not null;uniqueIndex"`
	Price     float64   `gorm:"type:numeric(12,2);not null"`
	Stock     int       `gorm:"column:stock;type:int;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
