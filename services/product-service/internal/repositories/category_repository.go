package repositories

import (
	"context"
	"encoding/json"
	"yaak-kaii/services/product-service/internal/models"
	"yaak-kaii/services/product-service/pkg/types"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	CreateTx(ctx context.Context, category *types.CreateCategoryPayload) (string, error)
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

func (o *categoryRepositoryImpl) CreateTx(
	ctx context.Context,
	category *types.CreateCategoryPayload,
) (string, error) {
	var categoryID string

	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		toCreateCat := &models.Category{
			Name:        category.Name,
			Description: category.Description,
		}

		if err := tx.Create(toCreateCat).Error; err != nil {
			return err
		}

		toCreateAttr := make([]models.CategoryAttribute, 0, len(category.Attributes))

		for _, item := range category.Attributes {
			optionsJSON, err := json.Marshal(item.Options)
			if err != nil {
				return err
			}

			toCreateAttr = append(toCreateAttr, models.CategoryAttribute{
				CategoryID: toCreateCat.ID,
				Key:        item.Key,
				Label:      item.Label,
				Type:       item.Type,
				Options:    datatypes.JSON(optionsJSON),
				Required:   item.Required,
				Scope:      item.Scope,
				AxisOrder:  int(item.AxisOrder),
			})
		}

		if len(toCreateAttr) > 0 {
			if err := tx.Create(&toCreateAttr).Error; err != nil {
				return err
			}
		}

		categoryID = toCreateCat.ID.String()
		return nil
	})

	if err != nil {
		return "", err
	}
	return categoryID, nil
}
