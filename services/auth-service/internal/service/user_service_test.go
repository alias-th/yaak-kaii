package service

import (
	"context"
	"errors"
	"testing"
	"yaak-kaii/services/auth-service/internal/domain"
	"yaak-kaii/services/auth-service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUser_Success(t *testing.T) {
	// 1. สร้าง Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. สร้าง Mock Repositories
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// 3. ตั้ง expectations ตามลำดับที่ service เรียก
	gomock.InOrder(
		// Step 1: ตรวจสอบว่า email มีแล้วหรือไม่
		mockUserRepo.EXPECT().
			GetUserByEmail(gomock.Any(), "john@example.com").
			Return(nil, nil). // ไม่มีแล้ว = nil, nil
			Times(1),

		// Step 2: ดึง role "user"
		mockRoleRepo.EXPECT().
			GetRoleByName(gomock.Any(), "user").
			Return(&domain.RoleModel{ID: uuid.New()}, nil).
			Times(1),

		// Step 3: สร้าง user
		mockUserRepo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(&domain.UserModel{
				ID:        uuid.New(),
				Email:     "john@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Role:      &domain.RoleModel{ID: uuid.New()},
			}, nil).
			Times(1),
	)

	// 4. สร้าง Service ด้วย mocks
	svc := &service{
		userRepo: mockUserRepo,
		roleRepo: mockRoleRepo,
	}

	// 5. เรียก function
	user, err := svc.CreateUser(context.Background(), &domain.CreateUserRequest{
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "0812345678",
		Password:    "password123",
	})

	// 6. ตรวจสอบผลลัพธ์
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.User.Email)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// ตั้ง expectations: GetUserByEmail คืน user ที่มีแล้ว
	mockUserRepo.EXPECT().
		GetUserByEmail(gomock.Any(), "existing@example.com").
		Return(&domain.UserModel{
			ID:    uuid.New(),
			Email: "existing@example.com",
		}, nil).
		Times(1)

	svc := &service{
		userRepo: mockUserRepo,
		roleRepo: mockRoleRepo,
	}

	user, err := svc.CreateUser(context.Background(), &domain.CreateUserRequest{
		Email:    "existing@example.com",
		Password: "password123",
	})

	// ต้องมี error
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user already exists", err.Error())
}

func TestCreateUser_RoleNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	gomock.InOrder(
		// Email ไม่มี
		mockUserRepo.EXPECT().
			GetUserByEmail(gomock.Any(), "newuser@example.com").
			Return(nil, nil).
			Times(1),

		// Role ไม่มี
		mockRoleRepo.EXPECT().
			GetRoleByName(gomock.Any(), "user").
			Return(nil, errors.New("role not found")).
			Times(1),
	)

	svc := &service{
		userRepo: mockUserRepo,
		roleRepo: mockRoleRepo,
	}

	user, err := svc.CreateUser(context.Background(), &domain.CreateUserRequest{
		Email:    "newuser@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, user)
}
