package repository

import (
	"context"
	"log"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
)

type roleRepository struct {
	store db.Store
}

func NewRoleRepository(store db.Store) domain.RoleRepository {
	return &roleRepository{store: store}
}

func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*domain.RoleModel, error) {
	role, err := r.store.GetRoleByName(ctx, name)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	return mapDBRoleToDomain(role), nil
}

func mapDBRoleToDomain(dbRole db.Role) *domain.RoleModel {
	return &domain.RoleModel{
		ID:          dbRole.ID,
		Name:        dbRole.Name,
		Description: dbRole.Description,
	}
}
