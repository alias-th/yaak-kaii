package service

import (
	"context"
	"yaak-kaii/services/product-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	productRepo domain.ProductRepository
}

func NewService(productRepo domain.ProductRepository) *service {
	return &service{
		productRepo: productRepo,
	}
}

func (s service) CreateProduct(ctx context.Context, product *domain.ProductModel) (*domain.ProductModel, error) {
	return s.productRepo.CreateProduct(ctx, product)
}

func (s service) GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.ProductModel, error) {
	return s.productRepo.GetProductByID(ctx, id)
}

func (s service) UpdateProduct(ctx context.Context, product *domain.ProductModel) (*domain.ProductModel, error) {
	return s.productRepo.UpdateProduct(ctx, product)
}

func (s service) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	return s.productRepo.DeleteProduct(ctx, id)
}

func (s service) ListProducts(ctx context.Context, filter map[string]interface{}) ([]*domain.ProductModel, error) {
	return s.productRepo.ListProducts(ctx, filter)
}
