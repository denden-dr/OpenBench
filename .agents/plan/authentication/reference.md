# Implementation Reference: Authentication Integration

This document provides focused reference code for implementing Supabase JWT verification and local user synchronization, as defined in `plan.md`. Application boilerplate is omitted.

## 1. `pkg/config/config.go` (Additions)

Add Supabase configurations to the `Config` struct and derive the JWKS endpoint:

```go
package config

import (
	"fmt"
	// ...
)

type Config struct {
	// ... database and pool fields ...
	SupabaseURL     string `envconfig:"SUPABASE_URL" required:"true"`
	SupabaseJWKSURL string `envconfig:"-"` // Calculated field
}

func LoadConfig() (*Config, error) {
	// ...
	
	// Derive JWKS URL
	cfg.SupabaseJWKSURL = fmt.Sprintf("%s/auth/v1/.well-known/jwks.json", cfg.SupabaseURL)

	return &cfg, nil
}
```

## 2. `internal/domain/user.go`

Define the domain model corresponding to the local `users` table:

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the local user profile synchronized from Supabase
type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	FullName  string    `db:"full_name" json:"full_name"`
	AvatarURL string    `db:"avatar_url" json:"avatar_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
```

## 3. `internal/repository/user_repo.go`

Set up the User repository to idempotently synchronize users:

```go
package repository

import (
	"context"
	"fmt"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	UpsertFromAuth(ctx context.Context, id uuid.UUID, email, fullName, avatarURL string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) UpsertFromAuth(ctx context.Context, id uuid.UUID, email, fullName, avatarURL string) (*domain.User, error) {
	query := `
		INSERT INTO users (id, email, full_name, avatar_url, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			full_name = EXCLUDED.full_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
		RETURNING *
	`
	var user domain.User
	err := r.db.QueryRowxContext(ctx, query, id, email, fullName, avatarURL).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	return &user, err
}
```

## 4. `internal/service/auth_service.go`

Service definition encapsulating JWT parsing, signature verification based on JWKS, and synching.

```go
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type AuthService interface {
	VerifyAndSync(ctx context.Context, rawToken string) (*domain.User, error)
}

type authService struct {
	keySet   jwk.Set
	userRepo repository.UserRepository
}

func NewAuthService(keySet jwk.Set, userRepo repository.UserRepository) AuthService {
	return &authService{
		keySet:   keySet,
		userRepo: userRepo,
	}
}

func (s *authService) VerifyAndSync(ctx context.Context, rawToken string) (*domain.User, error) {
	// Parse and verify signature/claims using the JWKS KeySet
	token, err := jwt.ParseString(rawToken, jwt.WithKeySet(s.keySet))
	if err != nil {
		return nil, fmt.Errorf("verifying token: %w", err)
	}

	// Extract subject (Supabase User UUID)
	sub, ok := token.Subject()
	if !ok {
		return nil, errors.New("subject claim missing")
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token: %w", err)
	}

	// Extract email claim
	var email string
	if err := token.Get("email", &email); err != nil {
		return nil, fmt.Errorf("email claim missing: %w", err)
	}

	// Extract optional metadata
	var fullName, avatarURL string
	var metadata map[string]interface{}
	if err := token.Get("user_metadata", &metadata); err == nil {
		if fn, ok := metadata["full_name"].(string); ok {
			fullName = fn
		}
		if av, ok := metadata["avatar_url"].(string); ok {
			avatarURL = av
		}
	}

	// Sync to local DB
	user, err := s.userRepo.UpsertFromAuth(ctx, userID, email, fullName, avatarURL)
	if err != nil {
		return nil, fmt.Errorf("syncing user: %w", err)
	}

	return user, nil
}
```

## 5. `internal/middleware/auth.go`

Fiber middleware to intercept protected calls. 

```go
package middleware

import (
	"strings"

	"github.com/denden-dr/OpenBench/internal/service"
	"github.com/gofiber/fiber/v3"
)

// RequireAuth returns a middleware that verifies the Supabase JWT
func RequireAuth(authService service.AuthService) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or malformed authorization header",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := authService.VerifyAndSync(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
```

## 6. `internal/handlers/auth.go`

Profile endpoint to verify local user retrieval.

```go
package handlers

import (
	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/gofiber/fiber/v3"
)

func GetMe(c fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not found in context",
		})
	}

	return c.JSON(user)
}
```

## 7. `cmd/api/main.go` (Wiring Updates)

Add these focused steps within the `main()` function:

```go
package main

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	// ... other imports ...
)

func main() {
	// ... config, logging, and database initializations ...

	// Initialize JWKS KeySet for Supabase Auth
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	keySet, err := jwk.Fetch(ctx, cfg.SupabaseJWKSURL)
	if err != nil {
		log.Fatal("Failed to fetch Supabase JWKS", zap.Error(err), zap.String("url", cfg.SupabaseJWKSURL))
	}
	log.Info("Successfully fetched Supabase JWKS")

	// Wire Auth Dependencies
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(keySet, userRepo)

	// ... fiber app instantiation ...

	// Protected Routes (v1)
	v1 := app.Group("/api/v1")
	{
		// Auth Routes (Protected)
		auth := v1.Group("/auth", middleware.RequireAuth(authService))
		auth.Get("/me", handlers.GetMe)
	}

	// ... server listen ...
}
```

## 8. `migrations/001_create_users_table.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    full_name TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```
