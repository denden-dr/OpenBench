package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/denden-dr/OpenBench/config"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockRepository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepository) BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error {
	args := m.Called(ctx, jti, expiresAt)
	return args.Error(0)
}

func (m *mockRepository) DeleteExpiredBlacklistedTokens(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}

var errDb = errors.New("db connection failure")

func TestAuthService_Login(t *testing.T) {
	password := "secretpassword123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	cfg := &config.Config{
		App: config.AppConfig{Env: "development"},
		Auth: config.AuthConfig{
			AccessSecret:  "access_secret_key_access_secret_key_access_secret_key",
			RefreshSecret: "refresh_secret_key_refresh_secret_key_refresh_secret_key",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 24 * time.Hour,
		},
	}

	user := &models.User{
		ID:           "u123",
		Email:        "admin@openbench.com",
		PasswordHash: string(hashed),
		FullName:     "Admin",
		Role:         "ADMIN",
	}

	tests := []struct {
		name        string
		email       string
		password    string
		mockErr     error
		expectedErr error
	}{
		{
			name:        "successful login",
			email:       "admin@openbench.com",
			password:    "secretpassword123",
			mockErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "wrong password",
			email:       "admin@openbench.com",
			password:    "wrongpassword",
			mockErr:     nil,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name:        "user not found",
			email:       "nonexistent@openbench.com",
			password:    "secretpassword123",
			mockErr:     nil,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name:        "empty email",
			email:       "",
			password:    "secretpassword123",
			mockErr:     nil,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name:        "empty password",
			email:       "admin@openbench.com",
			password:    "",
			mockErr:     nil,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name:        "repo error",
			email:       "admin@openbench.com",
			password:    "secretpassword123",
			mockErr:     errDb,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			switch {
			case tt.mockErr != nil:
				repo.On("GetUserByEmail", mock.Anything, tt.email).Return(nil, tt.mockErr)
			case tt.name == "user not found" || tt.name == "empty email":
				repo.On("GetUserByEmail", mock.Anything, tt.email).Return(nil, nil)
			default:
				repo.On("GetUserByEmail", mock.Anything, tt.email).Return(user, nil)
			}

			svc := NewService(repo, repo, cfg)

			result, err := svc.Login(context.Background(), tt.email, tt.password)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(result)
				is.NotEmpty(result.AccessToken)
				is.NotEmpty(result.RefreshToken)
				is.Equal(tt.email, result.User.Email)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	password := "secretpassword123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	cfg := &config.Config{
		App: config.AppConfig{Env: "development"},
		Auth: config.AuthConfig{
			AccessSecret:  "access_secret_key_access_secret_key_access_secret_key",
			RefreshSecret: "refresh_secret_key_refresh_secret_key_refresh_secret_key",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 24 * time.Hour,
		},
	}

	user := &models.User{
		ID:           "u123",
		Email:        "admin@openbench.com",
		PasswordHash: string(hashed),
		FullName:     "Admin",
		Role:         "ADMIN",
	}

	// We need a real token for refresh tests. We can mock repository to login successfully first.
	loginRepo := &mockRepository{}
	loginRepo.On("GetUserByEmail", mock.Anything, "admin@openbench.com").Return(user, nil)
	svc := NewService(loginRepo, loginRepo, cfg)

	loginResult, err := svc.Login(context.Background(), "admin@openbench.com", "secretpassword123")
	if err != nil {
		t.Fatalf("unexpected error during login: %v", err)
	}

	tests := []struct {
		name         string
		refreshToken string
		setupMock    func(repo *mockRepository)
		expectedErr  error
	}{
		{
			name:         "successful refresh",
			refreshToken: loginResult.RefreshToken,
			setupMock: func(repo *mockRepository) {
				repo.On("BlacklistToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil)
			},
			expectedErr: nil,
		},
		{
			name:         "invalid refresh token string",
			refreshToken: "invalid-token-string",
			setupMock:    func(repo *mockRepository) {},
			expectedErr:  ErrInvalidToken,
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			setupMock:    func(repo *mockRepository) {},
			expectedErr:  ErrInvalidToken,
		},
		{
			name:         "user not found in db during refresh",
			refreshToken: loginResult.RefreshToken,
			setupMock: func(repo *mockRepository) {
				repo.On("BlacklistToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("GetUserByEmail", mock.Anything, user.Email).Return(nil, nil)
			},
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := assert.New(t)
			must := require.New(t)

			repo := &mockRepository{}
			tt.setupMock(repo)
			testSvc := NewService(repo, repo, cfg)

			refreshResult, err := testSvc.Refresh(context.Background(), tt.refreshToken)
			if tt.expectedErr != nil {
				must.Error(err)
				is.ErrorIs(err, tt.expectedErr)
			} else {
				must.NoError(err)
				must.NotNil(refreshResult)
				is.NotEmpty(refreshResult.AccessToken)
			}
			repo.AssertExpectations(t)
		})
	}
}

func BenchmarkPasswordHashing(b *testing.B) {
	password := []byte("secretpassword123")
	cost := bcrypt.DefaultCost

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = bcrypt.GenerateFromPassword(password, cost)
	}
}

func BenchmarkPasswordComparison(b *testing.B) {
	password := []byte("secretpassword123")
	cost := bcrypt.DefaultCost
	hashed, _ := bcrypt.GenerateFromPassword(password, cost)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = bcrypt.CompareHashAndPassword(hashed, password)
	}
}
