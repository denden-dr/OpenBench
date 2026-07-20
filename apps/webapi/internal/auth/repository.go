package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/samber/hot"
)

type QueryRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
}

type CommandRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error
	DeleteExpiredBlacklistedTokens(ctx context.Context) (int64, error)
}

type sqlQueryRepository struct {
	db    *sqlx.DB
	cache *hot.HotCache[string, bool]
	psql  squirrel.StatementBuilderType
}

type sqlCommandRepository struct {
	db    *sqlx.DB
	cache *hot.HotCache[string, bool]
	psql  squirrel.StatementBuilderType
}

func NewQueryRepository(db *sqlx.DB, cache *hot.HotCache[string, bool]) QueryRepository {
	return &sqlQueryRepository{
		db:    db,
		cache: cache,
		psql:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func NewCommandRepository(db *sqlx.DB, cache *hot.HotCache[string, bool]) CommandRepository {
	return &sqlCommandRepository{
		db:    db,
		cache: cache,
		psql:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *sqlQueryRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query, args, err := r.psql.Select("id", "email", "password_hash", "full_name", "role", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(squirrel.Eq{"email": email, "deleted_at": nil}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *sqlQueryRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query, args, err := r.psql.Select("id", "email", "password_hash", "full_name", "role", "created_at", "updated_at", "deleted_at").
		From("users").
		Where(squirrel.Eq{"id": id, "deleted_at": nil}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *sqlCommandRepository) CreateUser(ctx context.Context, user *models.User) error {
	query, args, err := r.psql.Insert("users").
		Columns("id", "email", "password_hash", "full_name", "role", "created_at", "updated_at").
		Values(user.ID, user.Email, user.PasswordHash, user.FullName, user.Role, squirrel.Expr("NOW()"), squirrel.Expr("NOW()")).
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	_, err = querier.ExecContext(ctx, query, args...)
	return err
}

func (r *sqlQueryRepository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	if isBlacklisted, found, _ := r.cache.Get(jti); found {
		return isBlacklisted, nil
	}

	query, args, err := r.psql.Select("1").
		From("token_blacklists").
		Where(squirrel.Eq{"jti": jti}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.cache.Set(jti, false)
			return false, nil
		}
		return false, err
	}

	r.cache.Set(jti, true)
	return true, nil
}

func (r *sqlCommandRepository) BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error {
	query, args, err := r.psql.Insert("token_blacklists").
		Columns("jti", "expires_at").
		Values(jti, expiresAt).
		Suffix("ON CONFLICT (jti) DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	_, err = querier.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	r.cache.Set(jti, true)
	return nil
}

func (r *sqlCommandRepository) DeleteExpiredBlacklistedTokens(ctx context.Context) (int64, error) {
	query, args, err := r.psql.Delete("token_blacklists").
		Where(squirrel.LtOrEq{"expires_at": squirrel.Expr("NOW()")}).
		ToSql()
	if err != nil {
		return 0, err
	}

	querier := database.GetQuerier(ctx, r.db)
	result, err := querier.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}
