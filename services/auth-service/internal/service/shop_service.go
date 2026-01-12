package service

import (
	"context"
	"yaak-kaii/services/auth-service/internal/domain"

	"github.com/google/uuid"
)

func (s *service) GetShopUser(ctx context.Context, userId string) (*domain.ShopModel, error) {
	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	shop, err := s.shopRepo.GetShopByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return shop, nil
}
