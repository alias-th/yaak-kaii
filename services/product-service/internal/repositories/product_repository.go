package repositories

import (
	"context"
	"strconv"
	"strings"
	"yaak-kaii/services/product-service/internal/models"
	"yaak-kaii/services/product-service/internal/utils"
	"yaak-kaii/services/product-service/pkg/types"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProductTx(ctx context.Context, payload *CreateProductPayload) (string, error)
	ListProducts(ctx context.Context, query types.ListProductsReq) (*types.ListProductsRes, error)
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

type Category struct {
	CategoryID uuid.UUID
	Name       string
	Color      string
	Size       string
}
type CreateProductPayload struct {
	ShopID      uuid.UUID // ID of the shop creating the product
	UserID      string    // ID of the user creating the product
	Name        string    // Product name
	Description string    // Product description
	Price       float64   // Product price
	Stock       int32     // Available stock quantity
	Slug        string    // Product slug
	Attributes  types.Attributes
	Category    Category
}

func (o *productRepositoryImpl) CreateProductTx(ctx context.Context, payload *CreateProductPayload) (string, error) {
	productId := ""
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product := &models.Product{
			ShopID:      payload.ShopID,
			Name:        payload.Name,
			Description: payload.Description,
			Slug:        payload.Slug,
			CategoryID:  payload.Category.CategoryID,
		}
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		attByte, err := payload.Attributes.Value()
		if err != nil {
			return err
		}

		var nextNo int
		if err := tx.Raw(`
		SELECT COALESCE(pv.variant_no, 0) + 1
		FROM product_variants pv
		WHERE pv.product_id = ?;
		`, product.ID).Scan(&nextNo).Error; err != nil {
			return err
		}

		sku := utils.GenerateSKU(payload.Category.Name, payload.Category.Color, payload.Category.Size, nextNo)

		if err := tx.Create(&models.ProductVariant{
			ProductID:  product.ID,
			Sku:        sku,
			Price:      payload.Price,
			Stock:      int(payload.Stock),
			Attributes: attByte,
			VariantNo:  nextNo,
		}).Error; err != nil {
			return err
		}
		productId = product.ID.String()
		return nil
	})

	if err != nil {
		return "", err
	}
	return productId, nil
}

func (o *productRepositoryImpl) ListProducts(ctx context.Context, q types.ListProductsReq) (*types.ListProductsRes, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	offset := int((q.Page - 1) * q.Limit)

	var minPrice *string
	var maxPrice *string
	if q.MinPrice > 0 {
		s := strconv.FormatInt(q.MinPrice, 10)
		minPrice = &s
	}
	if q.MaxPrice > 0 {
		s := strconv.FormatInt(q.MaxPrice, 10)
		maxPrice = &s
	}

	base := o.db.WithContext(ctx).Model(&models.Product{})

	if s := strings.TrimSpace(q.Q); s != "" {
		base = base.Where("products.name ILIKE ?", "%"+s+"%")
	}

	if len(q.CategoryIds) > 0 {
		base = base.Where("products.category_id = ANY(?)", pq.Array(q.CategoryIds))
	}

	if minPrice != nil || maxPrice != nil {
		base = base.Where(`
			EXISTS (
				SELECT 1
				FROM product_variants pv
				WHERE pv.product_id = products.id
				  AND (?::numeric IS NULL OR pv.price >= ?::numeric)
				  AND (?::numeric IS NULL OR pv.price <= ?::numeric)
			)
		`, minPrice, minPrice, maxPrice, maxPrice)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var products []models.Product

	if err := base.Session(&gorm.Session{}).
		Order(utils.ProductOrder(q.Sort)).
		Limit(int(q.Limit)).
		Offset(offset).
		Preload("Category").
		Preload("Variants").
		Preload("Images", func(d *gorm.DB) *gorm.DB { return d.Order("created_at ASC") }).
		Find(&products).Error; err != nil {
		return nil, err
	}

	pagination := utils.NewPagination(int(q.Page), int(q.Limit), total)
	return &types.ListProductsRes{
		Products:   products,
		Pagination: *pagination,
	}, nil
}
