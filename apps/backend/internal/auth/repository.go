package auth

import (
	"context"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name=Repository --output=mocks --outpkg=mocks --case=underscore
type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetRefreshTokenWithLock(ctx context.Context, tx *sqlx.Tx, tokenHash string) (*RefreshTokenRecord, error)
	GetRefreshTokenByID(ctx context.Context, tx *sqlx.Tx, id string) (*RefreshTokenRecord, error)
	CreateRefreshToken(ctx context.Context, tx *sqlx.Tx, r *RefreshTokenRecord) error
	UpdateRefreshToken(ctx context.Context, tx *sqlx.Tx, r *RefreshTokenRecord) error
	RevokeTokenFamily(ctx context.Context, tx *sqlx.Tx, familyID string) error
	RevokeTokenByHash(ctx context.Context, tokenHash string) error
}

type postgresRepository struct {
	db *database.Database
}

func NewRepository(db *database.Database) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := "SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1"
	err := r.db.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	query := "SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1"
	err := r.db.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) GetRefreshTokenWithLock(ctx context.Context, tx *sqlx.Tx, tokenHash string) (*RefreshTokenRecord, error) {
	var record RefreshTokenRecord
	query := `
		SELECT id, family_id, user_id, token_hash, is_used, is_revoked, used_at, expires_at, created_at, replaced_by_token_id
		FROM refresh_tokens
		WHERE token_hash = $1
		FOR UPDATE
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &record, query, tokenHash)
	} else {
		err = r.db.DB.GetContext(ctx, &record, query, tokenHash)
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *postgresRepository) GetRefreshTokenByID(ctx context.Context, tx *sqlx.Tx, id string) (*RefreshTokenRecord, error) {
	var record RefreshTokenRecord
	query := `
		SELECT id, family_id, user_id, token_hash, is_used, is_revoked, used_at, expires_at, created_at, replaced_by_token_id
		FROM refresh_tokens
		WHERE id = $1
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &record, query, id)
	} else {
		err = r.db.DB.GetContext(ctx, &record, query, id)
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *postgresRepository) CreateRefreshToken(ctx context.Context, tx *sqlx.Tx, record *RefreshTokenRecord) error {
	query := `
		INSERT INTO refresh_tokens (id, family_id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, record.ID, record.FamilyID, record.UserID, record.TokenHash, record.ExpiresAt)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, record.ID, record.FamilyID, record.UserID, record.TokenHash, record.ExpiresAt)
	}
	return err
}

func (r *postgresRepository) UpdateRefreshToken(ctx context.Context, tx *sqlx.Tx, record *RefreshTokenRecord) error {
	query := `
		UPDATE refresh_tokens
		SET is_used = $1, used_at = $2, replaced_by_token_id = $3, is_revoked = $4
		WHERE id = $5
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, record.IsUsed, record.UsedAt, record.ReplacedByTokenID, record.IsRevoked, record.ID)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, record.IsUsed, record.UsedAt, record.ReplacedByTokenID, record.IsRevoked, record.ID)
	}
	return err
}

func (r *postgresRepository) RevokeTokenFamily(ctx context.Context, tx *sqlx.Tx, familyID string) error {
	query := "UPDATE refresh_tokens SET is_revoked = TRUE WHERE family_id = $1"
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, familyID)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, familyID)
	}
	return err
}

func (r *postgresRepository) RevokeTokenByHash(ctx context.Context, tokenHash string) error {
	query := "UPDATE refresh_tokens SET is_revoked = TRUE WHERE token_hash = $1"
	_, err := r.db.DB.ExecContext(ctx, query, tokenHash)
	return err
}
