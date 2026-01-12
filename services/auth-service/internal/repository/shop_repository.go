package repository

import (
	"context"
	"log"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"

	"github.com/google/uuid"
)

type shopRepository struct {
	store db.Store
}

func NewShopRepository(store db.Store) domain.ShopRepository {
	return &shopRepository{store: store}
}

func (r *shopRepository) GetShopByUserID(ctx context.Context, userID uuid.UUID) (*domain.ShopModel, error) {
	shop, err := r.store.GetShopByUserID(ctx, userID)
	if err != nil {
		log.Println("Error fetching shop by user ID:", err)
		return nil, err
	}
	return &domain.ShopModel{
		ID:          shop.ID.String(),
		SellerID:    shop.SellerID.String(),
		Name:        shop.Name,
		Slug:        shop.Slug,
		Description: shop.Description.String,
		IsActive:    shop.IsActive,
	}, nil
}
