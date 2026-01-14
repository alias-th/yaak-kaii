package repositories

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
	"yaak-kaii/services/product-service/internal/models"
	"yaak-kaii/services/product-service/internal/utils"
	"yaak-kaii/services/product-service/pkg/types"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProductTx(ctx context.Context, payload *CreateProductPayload) (string, error)
	CreateProductVariant(ctx context.Context, payload *types.CreateProductVariantPayload) (string, error)
	ListProducts(ctx context.Context, query types.ListProductsReq) (*types.ListProductsRes, error)
	GetProductByID(ctx context.Context, productID string) (*models.Product, error)
	GetProductBySlug(ctx context.Context, shopID string, slug string) (*types.ProductDetailResponse, error)
	AddProductImages(ctx context.Context, productID string, imageUrls []string) error
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
	Name        string    // Product name
	Description string    // Product description
	Slug        string    // Product slug
	CategoryID  uuid.UUID
}

func (o *productRepositoryImpl) CreateProductTx(ctx context.Context, payload *CreateProductPayload) (string, error) {
	productId := ""
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

// Get product by ID
func (o *productRepositoryImpl) GetProductByID(ctx context.Context, productID string) (*models.Product, error) {
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return nil, err
	}
	var product models.Product
	if err := o.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Images", func(d *gorm.DB) *gorm.DB { return d.Order("created_at ASC") }).
		First(&product, "id = ?", productUUID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (o *productRepositoryImpl) AddProductImages(ctx context.Context, productID string, imageUrls []string) error {
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return err
	}

	images := make([]models.ProductImage, 0, len(imageUrls))
	for _, url := range imageUrls {
		images = append(images, models.ProductImage{
			ProductID: productUUID,
			ImageURL:  url,
		})
	}

	if err := o.db.WithContext(ctx).Create(&images).Error; err != nil {
		return err
	}
	return nil
}

func (o *productRepositoryImpl) CreateProductVariant(ctx context.Context, payload *types.CreateProductVariantPayload) (string, error) {
	var variantID string
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// load product
		var p models.Product
		if err := tx.First(&p, "id = ?", payload.ProductID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("product not found")
			}
			return err
		}

		// load variant axes
		var attrs []models.CategoryAttribute
		if err := tx.Where("category_id = ? AND scope = ?", p.CategoryID, "VARIANT").Order("axis_order ASC").Find(&attrs).Error; err != nil {
			return err
		}
		if len(attrs) == 0 {
			return errors.New("category has no variant attributes")
		}
		axes := make([]types.AxisDef, 0, len(attrs))
		for _, a := range attrs {
			opts, err := utils.ParseOptionsAsSet(a.Options)
			if err != nil {
				return err
			}
			axes = append(axes, types.AxisDef{
				Key:      a.Key,
				Required: a.Required,
				Options:  opts,
				Order:    a.AxisOrder,
			})
		}
		sort.Slice(axes, func(i, j int) bool { return axes[i].Order < axes[j].Order })

		// check unknown attributes
		if err := utils.RejectUnknownKeys(payload.Attributes, axes); err != nil {
			return err
		}

		// build variant key
		varKey, err := utils.BuildVariantKeyAndValidate(payload.Attributes, axes)
		if err != nil {
			return err
		}

		// create sku
		var c models.Category
		if err := tx.First(&c, "id = ?", p.CategoryID).Error; err != nil {
			return err
		}
		var nextNo int
		if err := tx.Raw(`
		SELECT COALESCE(pv.variant_no, 0) + 1
		FROM product_variants pv
		WHERE pv.product_id = ?;
		`, p.ID).Scan(&nextNo).Error; err != nil {
			return err
		}
		sku := utils.GenerateSKU(c.Name, payload.Attributes["color"], payload.Attributes["size"], nextNo)

		// Convert map[string]string to bytes
		attByte, err := payload.Attributes.Value()
		if err != nil {
			return err
		}

		// create variant
		variant := &models.ProductVariant{
			ProductID:  p.ID,
			Sku:        sku,
			Price:      payload.Price,
			Stock:      int(payload.Stock),
			Attributes: attByte,
			VariantNo:  nextNo,
			VariantKey: varKey,
		}
		if err := tx.Create(variant).Error; err != nil {
			return err
		}

		variantID = variant.ID.String()

		return nil
	})

	if err != nil {
		return "", err
	}

	return variantID, nil
}

func (o *productRepositoryImpl) GetProductBySlug(ctx context.Context, shopID string, slug string) (*types.ProductDetailResponse, error) {
	shopUUID, err := uuid.Parse(shopID)
	if err != nil {
		return nil, err
	}

	// load product
	var p models.Product
	if err := o.db.WithContext(ctx).Where("shop_id = ? AND slug = ?", shopUUID, slug).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// load variant axes
	var catAttrs []models.CategoryAttribute
	if err := o.db.WithContext(ctx).
		Where("category_id = ? AND scope = ?", p.CategoryID, "VARIANT").
		Order("axis_order ASC").
		Find(&catAttrs).Error; err != nil {
		return nil, err
	}
	axes := make([]types.AxisDTO, 0, len(catAttrs))
	for _, a := range catAttrs {
		opts, err := utils.ParseStringArrayOptions(a.Options)
		if err != nil {
			return nil, err
		}
		axes = append(axes, types.AxisDTO{
			Key:      a.Key,
			Label:    a.Label,
			Type:     a.Type,
			Required: a.Required,
			Options:  opts,
			Order:    a.AxisOrder,
		})
	}

	// load variants
	var variants []models.ProductVariant
	if err := o.db.WithContext(ctx).Where("product_id = ?", p.ID).Order("created_at ASC").Find(&variants).Error; err != nil {
		return nil, err
	}
	var variantDTOs []types.VariantDTO
	for _, v := range variants {
		attrs, err := v.Attributes.Value()
		if err != nil {
			return nil, err
		}
		variantDTOs = append(variantDTOs, types.VariantDTO{
			ID:         v.ID,
			SKU:        v.Sku,
			Price:      v.Price,
			Stock:      int32(v.Stock),
			Attributes: attrs,
			VariantKey: v.VariantKey,
		})
	}

	// load images
	var images []models.ProductImage
	if err := o.db.WithContext(ctx).Where("product_id = ?", p.ID).Order("created_at ASC").Find(&images).Error; err != nil {
		return nil, err
	}
	imageUrls := make([]string, 0, len(images))
	for _, img := range images {
		imageUrls = append(imageUrls, img.ImageURL)
	}

	// build response
	res := &types.ProductDetailResponse{
		Product: struct {
			ID          uuid.UUID
			ShopID      uuid.UUID
			CategoryID  uuid.UUID
			Name        string
			Slug        string
			Description string
			Status      string
			CreatedAt   time.Time
		}{
			ID:          p.ID,
			ShopID:      p.ShopID,
			CategoryID:  p.CategoryID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt,
		},
		Axes:     axes,
		Variants: variantDTOs,
		Images:   imageUrls,
	}

	return res, nil
}
