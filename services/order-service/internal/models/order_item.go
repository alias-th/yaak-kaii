package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID         string
	ProductID       string
	Sku             string
	Quantity        int
	NameAtPurchase  string
	PriceAtPurchase float64
}
