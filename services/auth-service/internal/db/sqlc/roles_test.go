package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRoleByName(t *testing.T) {
	testCases := []struct {
		name         string
		expectedName string
	}{
		{
			name:         "Get User Role",
			expectedName: "user",
		},
	}

	for _, testCese := range testCases {
		ctx := context.Background()

		role, err := testStore.GetRoleByName(ctx, testCese.expectedName)
		require.NoError(t, err)
		require.NotEmpty(t, role)
		assert.Equal(t, testCese.expectedName, role.Name)
		assert.NotEqual(t, uuid.Nil, role.ID)
		assert.NotEmpty(t, role.Description)
	}
}
