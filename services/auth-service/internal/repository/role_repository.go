package repository

import (
	"context"
	"database/sql"
	"log"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
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
		if err == sql.ErrNoRows {
			return nil, utils.NewRoleNotFoundError()
		}
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
