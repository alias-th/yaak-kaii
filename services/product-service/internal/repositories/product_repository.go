package repositories

import (
	"context"
	"yaak-kaii/services/product-service/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (o *productRepositoryImpl) Create(ctx context.Context, product *models.Product) error {
	result := gorm.WithResult()
	err := gorm.G[models.Product](o.db, result).Create(ctx, product)
	return err
}
