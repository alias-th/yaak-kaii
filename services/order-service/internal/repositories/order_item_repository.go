package repositories

import (
	"context"
	"yaak-kaii/services/order-service/internal/models"

	"gorm.io/gorm"
)

type OrderItemRepository interface {
	Create(ctx context.Context, order *models.OrderItem) error
}

type orderItemImpl struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemImpl{
		db: db,
	}
}

func (o *orderItemImpl) Create(ctx context.Context, orderItem *models.OrderItem) error {
	result := gorm.WithResult()
	err := gorm.G[models.OrderItem](o.db, result).Create(ctx, orderItem)
	return err
}
