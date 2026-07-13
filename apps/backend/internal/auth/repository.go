package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type QueryRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
}

type CommandRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error
	DeleteExpiredBlacklistedTokens(ctx context.Context) error
}

type sqlQueryRepository struct {
	db *sqlx.DB
}

type sqlCommandRepository struct {
	db *sqlx.DB
}

func NewQueryRepository(db *sqlx.DB) QueryRepository {
	return &sqlQueryRepository{db: db}
}

func NewCommandRepository(db *sqlx.DB) CommandRepository {
	return &sqlCommandRepository{db: db}
}

func (r *sqlQueryRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *sqlCommandRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	querier := database.GetQuerier(ctx, r.db)
	_, err := querier.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.FullName, user.Role)
	return err
}

func (r *sqlQueryRepository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	query := `SELECT 1 FROM token_blacklists WHERE jti = $1 LIMIT 1`
	var exists int
	err := r.db.GetContext(ctx, &exists, query, jti)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *sqlCommandRepository) BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error {
	query := `INSERT INTO token_blacklists (jti, expires_at) VALUES ($1, $2) ON CONFLICT (jti) DO NOTHING`
	querier := database.GetQuerier(ctx, r.db)
	_, err := querier.ExecContext(ctx, query, jti, expiresAt)
	return err
}

func (r *sqlCommandRepository) DeleteExpiredBlacklistedTokens(ctx context.Context) error {
	query := `DELETE FROM token_blacklists WHERE expires_at <= NOW()`
	querier := database.GetQuerier(ctx, r.db)
	_, err := querier.ExecContext(ctx, query)
	return err
}
