package repository

import (
	"context"
	"encoding/json"
	"time"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

type guestRepository struct {
	store db.Store
}

func NewGuestRepository(store db.Store) domain.GuestRepository {
	return &guestRepository{store: store}
}

func (r *guestRepository) CreateGuest(ctx context.Context, guest *domain.GuestModel) (*domain.GuestModel, error) {

	metadataBytes, err := json.Marshal(guest.MetaData)
	if err != nil {
		return nil, err
	}

	arg := db.CreateGuestParams{
		TokenHash: guest.TokenHash,
		IpAddr:    guest.IpAddress,
		UserAgent: guest.UserAgent,
		Metadata:  metadataBytes,
		ExpiresAt: pgtype.Timestamptz{Time: time.Unix(guest.ExpiresAt, 0), Valid: true},
	}
	result, err := r.store.CreateGuest(ctx, arg)
	if err != nil {
		return nil, err
	}
	return mapDBGuestToDomain(result), nil
}

func (r *guestRepository) GetGuestByToken(ctx context.Context, token string) (*domain.GuestModel, error) {
	result, err := r.store.GetGuestByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return mapDBGuestToDomain(result), nil
}

func mapDBGuestToDomain(dbGuest db.Guest) *domain.GuestModel {
	return &domain.GuestModel{
		ID:        dbGuest.ID,
		TokenHash: dbGuest.TokenHash,
		IpAddress: dbGuest.IpAddr,
		UserAgent: dbGuest.UserAgent,
		CreatedAt: dbGuest.CreatedAt.Unix(),
		ExpiresAt: dbGuest.CreatedAt.Unix(),
	}
}
