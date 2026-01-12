package grpcserver

import (
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
