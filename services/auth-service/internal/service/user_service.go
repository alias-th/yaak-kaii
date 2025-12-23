package service

import (
	"context"
	"log"
	"time"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
)

func (s *service) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	// 1. Checking already exist email
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, utils.NewUserAlreadyExistsError()
	}

	if !utils.IsUserNotFoundError(err) {
		return nil, utils.NewInternalServerError()
	}

	// 2. Get role user
	role, err := s.roleRepo.GetRoleByName(ctx, "user")
	if err != nil {
		if utils.IsRoleNotFoundError(err) {
			return nil, err
		}
		return nil, utils.NewInternalServerError()
	}

	// 3. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.NewInternalServerError()
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
			return nil, err
		}
		return nil, utils.NewInternalServerError()
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
		if utils.IsUserNotFoundError(err) {
			return nil, utils.NewInvalidCredentialsError()
		}
		return nil, utils.NewInternalServerError()
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
		log.Printf("failed to generate token: error=%v", err)
		return nil, utils.NewInternalServerError()
	}

	// 4. Generate refresh token
	refreshToken, refreshTokenHashed, err := utils.GenerateTokenPair(32)
	if err != nil {
		log.Printf("failed to generate refresh token: error=%v", err)
		return nil, utils.NewInternalServerError()
	}

	// 5. Save refresh token hashed
	expiresAt := time.Now().Add(time.Hour * 24 * 30).Unix()
	arg := &domain.RefreshTokenModel{
		User: domain.UserModel{
			ID: user.ID,
		},
		TokenHash: refreshTokenHashed,
		ExpiresAt: expiresAt,
	}
	err = s.refreshTokenRepo.CreateRefreshToken(ctx, arg)
	if err != nil {
		log.Printf("failed to save refresh token: error=%v", err)
		return nil, utils.NewInternalServerError()
	}

	// 6. Return response
	return &domain.LoginResponse{
		UserID:       user.ID.String(),
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
