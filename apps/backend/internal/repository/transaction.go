package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	Commit() error
	Rollback() error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}
