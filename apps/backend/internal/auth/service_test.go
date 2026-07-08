package auth

import (
	"context"
	"testing"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type mockRepository struct {
	user *models.User
	err  error
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.user != nil && m.user.Email == email {
		return m.user, nil
	}
	return nil, nil
}

func (m *mockRepository) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}

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

	repo := &mockRepository{user: user}
	svc := NewService(repo, cfg)

	t.Run("successful login", func(t *testing.T) {
		result, err := svc.Login(context.Background(), "admin@openbench.com", "secretpassword123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.AccessToken == "" || result.RefreshToken == "" {
			t.Error("tokens should not be empty")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		_, err := svc.Login(context.Background(), "admin@openbench.com", "wrongpassword")
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := svc.Login(context.Background(), "nonexistent@openbench.com", "secretpassword123")
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})
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

	repo := &mockRepository{user: user}
	svc := NewService(repo, cfg)

	result, err := svc.Login(context.Background(), "admin@openbench.com", "secretpassword123")
	if err != nil {
		t.Fatalf("unexpected error during login: %v", err)
	}

	t.Run("successful refresh", func(t *testing.T) {
		refreshResult, err := svc.Refresh(context.Background(), result.RefreshToken)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if refreshResult.AccessToken == "" {
			t.Error("new access token should not be empty")
		}
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		_, err := svc.Refresh(context.Background(), "invalid-token-string")
		if err != ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken, got %v", err)
		}
	})
}
