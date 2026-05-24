package repository

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrDuplicate           = errors.New("resource already exists")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrDatabaseUnavailable = errors.New("database is unavailable")
)

// MapDatabaseError translates raw sql/pgx database errors into standard repository errors
func MapDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrDuplicate
		case "23503": // foreign_key_violation
			return ErrForeignKeyViolation
		case "57P01", "57P02", "57P03", "08000", "08003", "08006": // shutdown / connection issues
			return ErrDatabaseUnavailable
		}
	}

	return err
}
