package services

import "yaak-kaii/services/order-service/internal/repositories"

type OrderItemService struct {
	repo repositories.OrderItemRepository
}

func NewOrderItemService(repo repositories.OrderItemRepository) *OrderItemService {
	return &OrderItemService{
		repo: repo,
	}
}

func (s *OrderItemService) CreateOrderItem() {}
