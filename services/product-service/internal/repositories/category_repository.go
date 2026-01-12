package repositories

import (
	"context"
	"yaak-kaii/services/product-service/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
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

func (o *categoryRepositoryImpl) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	category, err := gorm.G[models.Category](o.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &models.Category{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}
