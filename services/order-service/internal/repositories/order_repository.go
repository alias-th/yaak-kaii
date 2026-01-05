package repositories

import (
	"context"
	"yaak-kaii/services/order-service/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{
		db: db,
	}
}

func (o *orderRepositoryImpl) Create(ctx context.Context, order *models.Order) error {
	result := gorm.WithResult()
	err := gorm.G[models.Order](o.db, result).Create(ctx, order)
	return err
}
