package service

import (
	"context"
	"yaak-kaii/services/auth-service/internal/domain"
)

func (s *service) GetRoleByName(ctx context.Context, name string) (*domain.RoleModel, error) {
	result, err := s.roleRepo.GetRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &domain.RoleModel{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}, nil
}
