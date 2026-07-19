package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

// DBQuerier defines the subset of sqlx.DB and sqlx.Tx methods that we use
// in our Command repositories to support transparent transaction propagation.
type DBQuerier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// InjectTx injects a transaction into the context.
func InjectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// ExtractTx extracts the transaction from the context.
func ExtractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

// GetQuerier returns the transaction from the context if it exists,
// otherwise it falls back to the database connection pool.
func GetQuerier(ctx context.Context, db *sqlx.DB) DBQuerier {
	if tx := ExtractTx(ctx); tx != nil {
		return tx
	}
	return db
}

// TxManager defines the contract for executing functions within a transaction.
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type sqlTxManager struct {
	db *sqlx.DB
}

// NewTxManager creates a new instance of TxManager.
func NewTxManager(db *sqlx.DB) TxManager {
	return &sqlTxManager{db: db}
}

// RunInTx executes a function inside a transaction.
// It handles rollback on error and commit on success.
// It is re-entrant: if a transaction is already active in the context,
// it will simply reuse it.
func (tm *sqlTxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx := ExtractTx(ctx); tx != nil {
		return fn(ctx)
	}

	tx, err := tm.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctxWithTx := InjectTx(ctx, tx)
	if err := fn(ctxWithTx); err != nil {
		return err
	}

	return tx.Commit()
}
