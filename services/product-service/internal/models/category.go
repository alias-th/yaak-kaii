package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Attributes  []CategoryAttribute `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}

type CategoryAttribute struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID uuid.UUID      `gorm:"column:category_id;type:uuid;not null"`
	Key        string         `gorm:"column:key;type:varchar(255);not null"`
	Label      string         `gorm:"column:label;type:varchar(255);not null"`
	Type       string         `gorm:"column:type;type:varchar(100);not null"`
	Required   bool           `gorm:"column:required;type:boolean;not null"`
	Options    datatypes.JSON `gorm:"column:options;type:jsonb"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
