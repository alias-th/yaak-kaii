package service

import (
	"context"
	"fmt"
	"log"
	"time"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
)

func (s *service) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	// 1. Checking already exist email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		return nil, utils.NewUserAlreadyExistsError()
	}

	// 2. Get role user
	role, err := s.roleRepo.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, utils.NewRoleNotFoundError()
	}

	// 3. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. Map domain
	arg := &domain.UserModel{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: hashedPassword,
		Role:         &domain.RoleModel{ID: role.ID},
	}

	// 5. Create user
	result, err := s.userRepo.CreateUser(ctx, arg)
	if err != nil {
		if utils.IsUserAlreadyExistsError(err) {
			return nil, utils.NewUserAlreadyExistsError()
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. Response
	return &domain.CreateUserResponse{
		User: &domain.UserModel{
			ID:            result.ID,
			Email:         result.Email,
			EmailVerified: result.EmailVerified,
			FirstName:     result.FirstName,
			LastName:      result.LastName,
			PhoneNumber:   result.PhoneNumber,
			IsActive:      result.IsActive,
			CreatedAt:     result.CreatedAt,
			UpdatedAt:     result.UpdatedAt,
			Role: &domain.RoleModel{
				ID: result.Role.ID,
			},
		},
	}, nil
}

func (s *service) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// 1. Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.NewInvalidCredentialsError()
	}
	if user == nil {
		return nil, utils.NewInvalidCredentialsError()
	}
	if !user.IsActive {
		return nil, utils.NewUserInactiveError()

	}

	// 2. Verify password
	err = utils.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		log.Printf("failed to checking password: error=%v", err)
		return nil, utils.NewInvalidCredentialsError()
	}

	// 3. Generate access token
	token, err := s.jwtAuth.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 4. Generate refresh token
	refreshToken, refreshTokenHashed, err := utils.GenerateTokenPair(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 5. Save refresh token hashed
	expiresAt := time.Now().Add(time.Hour * 24 * 30).Unix()
	arg := &domain.RefreshTokenModel{
		User:      db.User{ID: user.ID},
		TokenHash: refreshTokenHashed,
		ExpiresAt: expiresAt,
	}
	err = s.refreshTokenRepo.CreateRefreshToken(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// 6. Return response
	return &domain.LoginResponse{
		UserID:       user.ID.String(),
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
