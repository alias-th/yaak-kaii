package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"yaak-kaii/services/product-service/internal/models"
	"yaak-kaii/services/product-service/internal/utils"
	"yaak-kaii/services/product-service/pkg/types"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	CreateTx(ctx context.Context, category *types.CreateCategoryPayload) (string, error)
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
	GetAllCategories(ctx context.Context, query *types.ListCategoriesReq) (*types.ListCategoriesRes, error)
	GetCategoryID(ctx context.Context, categoryID string) (*models.Category, error)
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

func (p *categoryRepositoryImpl) GetAllCategories(ctx context.Context, q *types.ListCategoriesReq) (*types.ListCategoriesRes, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}

	offset := int((q.Page - 1) * q.Limit)
	base := p.db.WithContext(ctx).Model(&models.Category{})

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var categories []models.Category
	err := base.Session(&gorm.Session{}).WithContext(ctx).Model(&models.Category{}).
		Order("created_at ASC").
		Limit(int(q.Limit)).
		Offset(offset).
		Preload("Attributes", func(d *gorm.DB) *gorm.DB { return d.Order("axis_order ASC") }).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	pagination := utils.NewPagination(int(q.Page), int(q.Limit), total)

	return &types.ListCategoriesRes{
		Categories: categories,
		Pagination: *pagination,
	}, nil

}

func (p *categoryRepositoryImpl) GetCategoryID(ctx context.Context, categoryID string) (*models.Category, error) {
	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}

	var category models.Category
	if err := p.db.WithContext(ctx).Preload("Attributes").First(&category, "id = ?", categoryUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}
