package grpcserver

import (
	"context"
	"yaak-kaii/services/product-service/internal/services"
	pb "yaak-kaii/shared/proto/product"
)

type ProductGRPCServer struct {
	pb.UnimplementedProductServiceServer
	productService  *services.ProductService
	categoryService *services.CategoryService
}

func NewProductGRPCServer(product *services.ProductService, category *services.CategoryService) *ProductGRPCServer {
	return &ProductGRPCServer{
		productService:  product,
		categoryService: category,
	}
}

func (s *ProductGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	payload := &services.CreateProductPayload{
		Name:        req.GetName(),
		UserId:      req.GetUserId(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		CategoryId:  req.GetCategoryId(),
		Stock:       req.GetStock(),
		Attributes:  req.GetAttributes(),
	}

	category, err := s.categoryService.GetCategoryByID(ctx, req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	categoryName := category.Name
	payload.CategoryName = categoryName

	err = s.productService.CreateProduct(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductResponse{
		ProductId: "new-product-id",
	}, nil
}
