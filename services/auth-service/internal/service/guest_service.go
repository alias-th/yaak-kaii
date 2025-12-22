package service

import (
	"context"
	"net/netip"
	"time"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"

	"github.com/ua-parser/uap-go/uaparser"
)

func (s *service) CreateGuest(ctx context.Context, guest *domain.CreateGuestRequest) (string, error) {
	// 1. Generate token
	token, tokenHashed, err := utils.GenerateTokenPair(32)
	if err != nil {
		return "", err
	}

	// 2. Parse IP address string into netip.Addr
	ip, err := netip.ParseAddr(guest.IpAddress)
	if err != nil {
		return "", err
	}

	// 3. Parse user-agent
	parser, err := uaparser.New()
	if err != nil {
		return "", err
	}
	client := parser.Parse(guest.UserAgent)

	// 4. Map to params
	now := time.Now().UTC()
	metaData := domain.MetaDataGuest{
		Os:               client.Os.ToString(),
		UserAgent:        client.UserAgent.ToString(),
		Device:           client.Device.ToString(),
		LastActivityTime: now.Unix(),
		SessionStartTime: now.Unix(),
		RequestCount:     1,
	}

	arg := &domain.GuestModel{
		TokenHash: tokenHashed,
		IpAddress: ip,
		UserAgent: guest.UserAgent,
		MetaData:  metaData,
		ExpiresAt: now.AddDate(0, 0, 1).Unix(),
	}

	_, err = s.guestRepo.CreateGuest(ctx, arg)
	if err != nil {
		return "", err
	}

	// 5. response
	return token, nil
}

func (s *service) VerifyGuestToken(ctx context.Context, token string) (*domain.GuestModel, error) {
	// decode token
	raw, err := utils.DecodeBase64URL(token)
	if err != nil {
		return nil, err
	}

	// hash token
	hashed := utils.HashTokenBytes(raw)

	guest, err := s.guestRepo.GetGuestByToken(ctx, hashed)
	if err != nil {
		return nil, err
	}

	// verify token
	ok, err := utils.VerifyToken(token, guest.TokenHash)
	if err != nil || !ok {
		return nil, err
	}

	return &domain.GuestModel{
		ID:        guest.ID,
		IpAddress: guest.IpAddress,
		UserAgent: guest.UserAgent,
		CreatedAt: guest.CreatedAt,
		ExpiresAt: guest.CreatedAt,
	}, nil

}
