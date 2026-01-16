package services

import (
	"context"
	"yaak-kaii/services/product-service/internal/models"
	"yaak-kaii/services/product-service/internal/repositories"
	"yaak-kaii/services/product-service/pkg/types"
)

type CategoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, payload *types.CreateCategoryPayload) (string, error) {
	id, err := s.repo.CreateTx(ctx, payload)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}
