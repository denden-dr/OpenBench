//go:build integration

package auth_test

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRepository_Integration(t *testing.T) {
	ctx := context.Background()

	// Spin up test db container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	cmdRepo := auth.NewCommandRepository(db)
	queryRepo := auth.NewQueryRepository(db)

	t.Run("CreateUser and GetUserByEmail", func(t *testing.T) {
		err := testutils.CleanTable(db, "users")
		require.NoError(t, err)

		user := &models.User{
			ID:           uuid.New().String(),
			Email:        "integration@test.com",
			PasswordHash: "hashedpassword",
			FullName:     "Integration Test User",
			Role:         "user",
		}

		err = cmdRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		fetched, err := queryRepo.GetUserByEmail(ctx, user.Email)
		require.NoError(t, err)
		require.NotNil(t, fetched)

		assert.Equal(t, user.ID, fetched.ID)
		assert.Equal(t, user.Email, fetched.Email)
		assert.Equal(t, user.PasswordHash, fetched.PasswordHash)
		assert.Equal(t, user.FullName, fetched.FullName)
		assert.Equal(t, user.Role, fetched.Role)
		assert.NotEmpty(t, fetched.CreatedAt)
		assert.NotEmpty(t, fetched.UpdatedAt)
		assert.Nil(t, fetched.DeletedAt)
	})

	t.Run("GetUserByEmail - Not Found", func(t *testing.T) {
		err := testutils.CleanTable(db, "users")
		require.NoError(t, err)

		fetched, err := queryRepo.GetUserByEmail(ctx, "nonexistent@test.com")
		require.NoError(t, err)
		assert.Nil(t, fetched)
	})
}
