package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Slug          string             `bson:"slug"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	Category      string             `bson:"category"`
	Brand         string             `bson:"brand"`
	StockQuantity int64              `bson:"stockQuantity"`
	Attributes    Attributes         `bson:"attributes"`
	Variants      Variants           `bson:"variants"`
	Images        []string           `bson:"images"`
	CreatedAt     string             `bson:"createdAt"`
	UpdatedAt     string             `bson:"updatedAt"`
}

type Attributes struct {
	Material string `bson:"material"`
	Gender   string `bson:"gender"`
	Pattern  string `bson:"pattern"`
}

type Variants struct {
	Color string `bson:"color"`
	Size  string `bson:"size"`
	SKU   string `bson:"sku"`
	Stock int64  `bson:"stock"`
}
