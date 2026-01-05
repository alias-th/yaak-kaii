package models

import (
	"time"

	"gorm.io/gorm"
)

// Order Has many OrderItems, OrderID is the foreign key in OrderItem
type Order struct {
	gorm.Model
	OrderItems      []OrderItem
	UserID          string
	GuestID         string
	Status          string
	ShippingAddress ShippingAddress `gorm:"type:jsonb"`
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

type ShippingAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}
