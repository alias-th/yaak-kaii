package service

import (
	"context"
	"errors"
	"log"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/shared/utils"
)

func (s *service) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	// 1. Checking already exist email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	if user != nil {
		return nil, errors.New("user already exists")
	}

	// 2. Get role user
	role, err := s.roleRepo.GetRoleByName(ctx, "user")
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	if role == nil {
		return nil, errors.New("user role not found")
	}

	// 3. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
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
		log.Printf("Error: %v", err)
		return nil, err
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
	// TODO: Implement login logic
	// 1. Get user by email
	// 2. Verify password
	// 3. Generate access token
	// 4. Generate refresh token
	// 5. Return response
	return nil, errors.New("not implemented")
}
