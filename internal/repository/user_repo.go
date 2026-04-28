package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/denden-dr/OpenBench/internal/domain"
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
