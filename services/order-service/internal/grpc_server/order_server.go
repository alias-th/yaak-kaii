package grpcserver

import (
	"context"
	"yaak-kaii/services/order-service/internal/services"
	pb "yaak-kaii/shared/proto/order"
)

type OrderGRPCServer struct {
	pb.UnimplementedOrderServiceServer
	orderService     *services.OrderService
	orderItemService *services.OrderItemService
}

func NewOrderGRPCServer(order *services.OrderService, orderItem *services.OrderItemService) *OrderGRPCServer {
	return &OrderGRPCServer{
		orderService:     order,
		orderItemService: orderItem,
	}
}

func (o *OrderGRPCServer) CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	return &pb.CreateOrderResponse{
		OrderId:     "",
		UserId:      "",
		GuestId:     "",
		TotalAmount: "",
		Currency:    "",
		Status:      "",
		ShippingAddress: &pb.ShippingAddress{
			Street:     "",
			City:       "",
			State:      "",
			PostalCode: "",
			Country:    "",
		},
		CreatedAt: 0,
		UpdatedAt: 0,
	}, nil
}
