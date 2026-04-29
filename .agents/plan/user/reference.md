# User Local Database Reference Implementation

This document contains the template code to fulfill the `plan.md` requirements for the User local database integration.

## 1. Database Migration
**File**: `migrations/001_create_users_table.sql`
*(Exact numbering based on your migration setup)*

```sql
-- Up Migration
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    full_name TEXT,
    avatar_url TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for email lookups (optional but recommended if querying by email)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Down Migration
-- DROP INDEX IF EXISTS idx_users_email;
-- DROP TABLE IF EXISTS users;
```

---

## 2. Domain Model
**File**: `internal/domain/user.go`

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the local proxy of a Supabase authenticated user.
type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	FullName  *string   `db:"full_name" json:"full_name,omitempty"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
```

---

## 3. Repository Implementation
**File**: `internal/repository/user_repo.go`

This code strictly follows the constructor injection and public-interface/private-struct rules.

```go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	
	// Adjust this import to your actual project module path
	"github.com/yourusername/openbench/internal/domain"
)

// ErrUserNotFound is returned when a requested user does not exist in the local db.
var ErrUserNotFound = errors.New("user not found")

// UserRepository defines the data access contract for users.
type UserRepository interface {
	// UpsertFromAuth atomically creates or updates a user profile from Supabase claims.
	UpsertFromAuth(ctx context.Context, id uuid.UUID, email string, fullName, avatarURL *string) (*domain.User, error)
	// FindByID retrieves a single user by their Supabase UUID.
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

// userRepository is the private implementation of UserRepository.
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository constructs a new UserRepository using the provided database connection.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// UpsertFromAuth creates a new user or updates an existing one based on the Supabase ID.
func (r *userRepository) UpsertFromAuth(ctx context.Context, id uuid.UUID, email string, fullName, avatarURL *string) (*domain.User, error) {
	query := `
		INSERT INTO users (id, email, full_name, avatar_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			full_name = EXCLUDED.full_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
		RETURNING id, email, full_name, avatar_url, updated_at
	`

	var user domain.User
	// StructScan maps the RETURNING fields directly into our domain.User struct based on the `db` tags
	err := r.db.QueryRowxContext(ctx, query, id, email, fullName, avatarURL).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("upserting user from auth: %w", err)
	}

	return &user, nil
}

// FindByID retrieves a user by their ID. Returns ErrUserNotFound if not present.
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, full_name, avatar_url, updated_at 
		FROM users 
		WHERE id = $1
	`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("finding user by id: %w", err)
	}

	return &user, nil
}
```

---

## 4. Service Implementation
**File**: `internal/service/user_service.go`

```go
package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/repository"
)

// UserService defines the business logic for users.
type UserService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user profile: %w", err)
	}
	return user, nil
}
```

---

## 5. Handler Implementation
**File**: `internal/handlers/user_handler.go`

```go
package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/denden-dr/OpenBench/internal/service"
)

// UserHandler handles HTTP requests for users.
type UserHandler interface {
	GetMe(c fiber.Ctx) error
}

type userHandler struct {
	service service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) GetMe(c fiber.Ctx) error {
	// Assuming the user ID is stored in the context by AuthMiddleware
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	user, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
```

---

## 6. Service Unit Testing (with Mocking)
**File**: `internal/service/user_service_test.go`

```go
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/denden-dr/OpenBench/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpsertFromAuth(ctx context.Context, id uuid.UUID, email string, fullName, avatarURL *string) (*domain.User, error) {
	args := m.Called(ctx, id, email, fullName, avatarURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestUserService_GetProfile(t *testing.T) {
	userID := uuid.New()
	mockUser := &domain.User{ID: userID, Email: "test@example.com"}
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func(m *MockUserRepository)
		expectedUser  *domain.User
		expectedError string
	}{
		{
			name:   "Success",
			userID: userID,
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(mockUser, nil)
			},
			expectedUser:  mockUser,
			expectedError: "",
		},
		{
			name:   "User Not Found",
			userID: userID,
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, repository.ErrUserNotFound)
			},
			expectedUser:  nil,
			expectedError: repository.ErrUserNotFound.Error(),
		},
		{
			name:   "Repository Failure",
			userID: userID,
			mockSetup: func(m *MockUserRepository) {
				m.On("FindByID", ctx, userID).Return(nil, errors.New("db failure"))
			},
			expectedUser:  nil,
			expectedError: "getting user profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepository)
			tt.mockSetup(repo)

			svc := NewUserService(repo)
			user, err := svc.GetProfile(ctx, tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
			repo.AssertExpectations(t)
		})
	}
}
```

---

## 7. Repository Unit Testing (with sqlmock)
**File**: `internal/repository/user_repo_test.go`

```go
package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func()
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:   "Success",
			userID: userID,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "full_name", "avatar_url", "updated_at"}).
					AddRow(userID, "test@example.com", "Test User", nil, time.Now())

				mock.ExpectQuery("SELECT id, email, full_name, avatar_url, updated_at FROM users WHERE id = \\$1").
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedUser: &domain.User{ID: userID, Email: "test@example.com"},
			expectedError: nil,
		},
		{
			name:   "Not Found",
			userID: userID,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, email, full_name, avatar_url, updated_at FROM users WHERE id = \\$1").
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := repo.FindByID(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
			}
		})
	}
}

func TestUserRepository_UpsertFromAuth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()
	userID := uuid.New()
	email := "test@example.com"
	fullName := "Test User"

	tests := []struct {
		name          string
		userID        uuid.UUID
		email         string
		fullName      *string
		mockSetup     func()
		expectedUser  *domain.User
	}{
		{
			name:     "Success",
			userID:   userID,
			email:    email,
			fullName: &fullName,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "full_name", "avatar_url", "updated_at"}).
					AddRow(userID, email, fullName, nil, time.Now())

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(userID, email, &fullName, nil).
					WillReturnRows(rows)
			},
			expectedUser: &domain.User{ID: userID, Email: email},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := repo.UpsertFromAuth(ctx, tt.userID, tt.email, tt.fullName, nil)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedUser.ID, user.ID)
		})
	}
}
```
