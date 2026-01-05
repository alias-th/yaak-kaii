package services

import (
	"context"
	"yaak-kaii/services/order-service/internal/models"
	"yaak-kaii/services/order-service/internal/repositories"
)

type OrderService struct {
	repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	payload := models.Order{}
	err := s.repo.Create(ctx, &payload)
	return err
}
