package repository

import (
	"context"
	"yaak-kaii/services/product-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *mongoRepository {
	return &mongoRepository{
		db: db,
	}
}

func (mongo mongoRepository) CreateProduct(ctx context.Context, product *domain.ProductModel) (*domain.ProductModel, error) {
	return nil, nil
}

func (mongo mongoRepository) GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.ProductModel, error) {
	return nil, nil
}

func (mongo mongoRepository) UpdateProduct(ctx context.Context, product *domain.ProductModel) (*domain.ProductModel, error) {
	return nil, nil
}

func (mongo mongoRepository) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (mongo mongoRepository) ListProducts(ctx context.Context, filter map[string]interface{}) ([]*domain.ProductModel, error) {
	return nil, nil
}
