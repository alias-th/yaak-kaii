package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *ProductModel) (*ProductModel, error)
	GetProductByID(ctx context.Context, id primitive.ObjectID) (*ProductModel, error)
	UpdateProduct(ctx context.Context, product *ProductModel) (*ProductModel, error)
	DeleteProduct(ctx context.Context, id primitive.ObjectID) error
	ListProducts(ctx context.Context, filter map[string]interface{}) ([]*ProductModel, error)
}
