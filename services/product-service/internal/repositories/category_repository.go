package repositories

import (
	"context"
	"yaak-kaii/services/product-service/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (o *categoryRepositoryImpl) Create(ctx context.Context, category *models.Category) error {
	result := gorm.WithResult()
	err := gorm.G[models.Category](o.db, result).Create(ctx, category)
	return err
}
