package grpc

import (
	"context"
	"yaak-kaii/services/product-service/internal/domain"
	pb "yaak-kaii/shared/proto/product"

	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedProductServiceServer
	service domain.ProductService
}

func NewGRPCHandler(server *grpc.Server, service domain.ProductService) {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterProductServiceServer(server, handler)
}

func (h *gRPCHandler) CreateProduct(ctx context.Context,
	req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	// Implement the logic to create a product using h.service
	return nil, nil
}
