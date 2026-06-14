package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/denden-dr/openbench/apps/backend/internal/auth"
	"github.com/denden-dr/openbench/apps/backend/internal/auth/mocks"
	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupServiceTest(t *testing.T) (*mocks.Repository, auth.Service, sqlmock.Sqlmock) {
	mockDB, mockSQL, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "postgres")
	dbWrapper := &database.Database{DB: sqlxDB}

	repo := mocks.NewRepository(t)
	jwtSecret := "my_test_secret_key"
	service := auth.NewService(repo, dbWrapper, jwtSecret)

	t.Cleanup(func() {
		mockDB.Close()
		assert.NoError(t, mockSQL.ExpectationsWereMet())
	})

	return repo, service, mockSQL
}

func TestService_SignIn(t *testing.T) {
	repo, service, _ := setupServiceTest(t)
	ctx := context.Background()

	rawPassword := "SecurePassword123!"
	hashedPassword, err := auth.HashPassword(rawPassword)
	require.NoError(t, err)

	testUser := &auth.User{
		ID:           "user-1",
		Email:        "user@openbench.dev",
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("GetUserByEmail", ctx, "user@openbench.dev").Return(testUser, nil).Once()
		repo.On("CreateRefreshToken", ctx, mock.Anything, mock.Anything).Return(nil).Once()

		res, err := service.SignIn(ctx, "user@openbench.dev", rawPassword, 5*time.Minute, 24*time.Hour)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "user-1", res.User.ID)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RawRefreshToken)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		repo.On("GetUserByEmail", ctx, "wrong@openbench.dev").Return(nil, errors.New("not found")).Once()

		res, err := service.SignIn(ctx, "wrong@openbench.dev", rawPassword, 5*time.Minute, 24*time.Hour)
		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
		assert.Nil(t, res)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		repo.On("GetUserByEmail", ctx, "user@openbench.dev").Return(testUser, nil).Once()

		res, err := service.SignIn(ctx, "user@openbench.dev", "WrongPassword", 5*time.Minute, 24*time.Hour)
		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
		assert.Nil(t, res)
	})
}

func TestService_SignOut(t *testing.T) {
	repo, service, _ := setupServiceTest(t)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		token := "my_refresh_token"
		hashed := auth.HashSha256(token)
		repo.On("RevokeTokenByHash", ctx, hashed).Return(nil).Once()

		err := service.SignOut(ctx, token)
		assert.NoError(t, err)
	})

	t.Run("EmptyToken", func(t *testing.T) {
		err := service.SignOut(ctx, "")
		assert.NoError(t, err)
	})
}

func TestService_GetUserByID(t *testing.T) {
	repo, service, _ := setupServiceTest(t)
	ctx := context.Background()

	testUser := &auth.User{
		ID:    "user-123",
		Email: "user@openbench.dev",
		Role:  "user",
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("GetUserByID", ctx, "user-123").Return(testUser, nil).Once()

		u, err := service.GetUserByID(ctx, "user-123")
		require.NoError(t, err)
		assert.Equal(t, testUser, u)
	})

	t.Run("Error", func(t *testing.T) {
		repo.On("GetUserByID", ctx, "user-123").Return(nil, errors.New("db error")).Once()

		u, err := service.GetUserByID(ctx, "user-123")
		assert.Error(t, err)
		assert.Nil(t, u)
	})
}
