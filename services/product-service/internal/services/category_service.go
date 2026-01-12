package services

import (
	"context"
	"yaak-kaii/services/product-service/internal/models"
	"yaak-kaii/services/product-service/internal/repositories"
)

type CategoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) CreateCategory() {}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}
