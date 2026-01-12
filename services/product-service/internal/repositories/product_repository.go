package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"yaak-kaii/services/product-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProductTx(ctx context.Context, payload *CreateProductPayload) error
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

type CreateProductPayload struct {
	ShopID      uuid.UUID // ID of the shop creating the product
	UserID      string    // ID of the user creating the product
	Name        string    // Product name
	Description string    // Product description
	Price       float64   // Product price
	CategoryID  uuid.UUID // Product category
	Stock       int32     // Available stock quantity
	Slug        string    // Product slug
	Sku         string    // Product SKU
	Attributes  Attributes
}

type Attributes map[string]string

func (a Attributes) Value() ([]byte, error) {
	return json.Marshal(a)
}
func (a *Attributes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
func (o *productRepositoryImpl) CreateProductTx(ctx context.Context, payload *CreateProductPayload) error {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product := &models.Product{
			ShopID:      payload.ShopID,
			Name:        payload.Name,
			Description: payload.Description,
			Slug:        payload.Slug,
			CategoryID:  payload.CategoryID,
		}
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		value, err := payload.Attributes.Value()
		if err != nil {
			return err
		}

		if err := tx.Create(&models.ProductVariant{
			ProductID:  product.ID,
			Sku:        payload.Sku,
			Price:      payload.Price,
			Stock:      int(payload.Stock),
			Attributes: value,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
