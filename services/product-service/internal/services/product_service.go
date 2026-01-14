package services

import (
	"context"
	grpcclients "yaak-kaii/services/product-service/internal/grpc_clients"
	"yaak-kaii/services/product-service/internal/models"
	"yaak-kaii/services/product-service/internal/repositories"
	"yaak-kaii/services/product-service/pkg/types"
	"yaak-kaii/shared/proto/auth"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type ProductService struct {
	repo        repositories.ProductRepository
	grpcClients *grpcclients.GrpcClients
}

func NewProductService(repo repositories.ProductRepository, grpcClients *grpcclients.GrpcClients) *ProductService {
	return &ProductService{repo: repo, grpcClients: grpcClients}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *types.CreateProductPayload) (string, error) {
	shop, err := s.grpcClients.Auth.Client.GetShopUser(ctx, &auth.GetShopRequest{UserId: product.UserId})
	if err != nil {
		return "", err
	}
	shopID, err := uuid.Parse(shop.Id)
	if err != nil {
		return "", err
	}
	categoryID, err := uuid.Parse(product.CategoryId)
	if err != nil {
		return "", err
	}

	slugValue := slug.Make(product.Name)

	productId, err := s.repo.CreateProductTx(ctx, &repositories.CreateProductPayload{
		ShopID:      shopID,
		Name:        product.Name,
		Description: product.Description,
		Slug:        slugValue,
		CategoryID:  categoryID,
	})
	if err != nil {
		return "", err
	}
	return productId, nil
}

func (s *ProductService) ListSellerProducts(ctx context.Context, query types.ListSellerProductsReq) (*types.ListProductsRes, error) {
	shop, err := s.grpcClients.Auth.Client.GetShopUser(ctx, &auth.GetShopRequest{UserId: query.UserID})
	if err != nil {
		return nil, err
	}
	shopID, err := uuid.Parse(shop.Id)
	if err != nil {
		return nil, err
	}

	products, err := s.repo.ListSellerProducts(ctx, repositories.ListSellerProductsReq{
		ShopID: shopID,
		Limit:  query.Limit,
		Page:   query.Page,
		Sort:   query.Sort,
	})

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) ListBuyerProducts(ctx context.Context, query types.ListProductsReq) (*types.ListProductsRes, error) {
	products, err := s.repo.ListBuyerProducts(ctx, query)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, productID string) (*models.Product, error) {
	products, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductService) UploadProductImages(ctx context.Context, productID string, imageUrls []string) error {
	err := s.repo.AddProductImages(ctx, productID, imageUrls)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) CreateProductVariant(ctx context.Context, payload *types.CreateProductVariantPayload) (string, error) {
	variantID, err := s.repo.CreateProductVariant(ctx, &types.CreateProductVariantPayload{
		ProductID:  payload.ProductID,
		Price:      payload.Price,
		Stock:      payload.Stock,
		Attributes: payload.Attributes,
	})
	if err != nil {
		return "", err
	}

	return variantID, nil
}

func (s *ProductService) GetProductBySlug(ctx context.Context, shopID string, slug string) (*types.ProductDetailResponse, error) {
	productDetail, err := s.repo.GetProductBySlug(ctx, shopID, slug)
	if err != nil {
		return nil, err
	}

	return productDetail, nil
}
