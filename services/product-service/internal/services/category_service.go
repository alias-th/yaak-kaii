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

func (s *CategoryService) ListCategories(ctx context.Context, q *types.ListCategoriesReq) (*types.ListCategoriesRes, error) {
	categories, err := s.repo.GetAllCategories(ctx, q)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) GetCategoryID(ctx context.Context, categoryID string) (*models.Category, error) {
	category, err := s.repo.GetCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	return category, nil
}
