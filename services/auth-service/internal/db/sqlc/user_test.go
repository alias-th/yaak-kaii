package db

import (
	"context"
	"testing"
	"time"
	"yaak-kaii/shared/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	passwordHash, err := utils.HashPassword(utils.RandomString(8))
	require.NoError(t, err)
	ctx := context.Background()

	role, err := testStore.GetRoleByName(ctx, "user")
	require.NoError(t, err)
	require.NotEmpty(t, role)
	assert.Equal(t, "user", role.Name)

	arg := CreateUserParams{
		Email:        utils.RandomEmail(),
		PasswordHash: passwordHash,
		FirstName:    utils.RandomString(6),
		LastName:     utils.RandomString(6),
		PhoneNumber:  utils.RandomString(10),
		RoleID:       role.ID,
	}

	user, err := testStore.CreateUser(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	// data validation
	assert.Equal(t, arg.Email, user.Email)
	assert.Equal(t, arg.PasswordHash, user.PasswordHash)
	assert.Equal(t, arg.FirstName, user.FirstName)
	assert.Equal(t, arg.LastName, user.LastName)
	assert.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	assert.Equal(t, role.ID, user.RoleID)

	// default
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.True(t, user.IsActive)
	assert.False(t, user.EmailVerified)

	// time
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
	assert.False(t, user.DeletedAt.Valid)

	return user
}
func TestCreateAccount(t *testing.T) {
	createRandomUser(t)
}
