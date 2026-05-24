package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestMapDatabaseError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			name:     "nil error",
			input:    nil,
			expected: nil,
		},
		{
			name:     "sql.ErrNoRows",
			input:    sql.ErrNoRows,
			expected: ErrNotFound,
		},
		{
			name:     "unique violation PgError",
			input:    &pgconn.PgError{Code: "23505"},
			expected: ErrDuplicate,
		},
		{
			name:     "foreign key violation PgError",
			input:    &pgconn.PgError{Code: "23503"},
			expected: ErrForeignKeyViolation,
		},
		{
			name:     "database unavailable pgError",
			input:    &pgconn.PgError{Code: "08006"},
			expected: ErrDatabaseUnavailable,
		},
		{
			name:     "generic unknown error",
			input:    errors.New("some DB issue"),
			expected: errors.New("some DB issue"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapDatabaseError(tt.input)
			if tt.expected == nil {
				assert.NoError(t, result)
			} else {
				assert.Equal(t, tt.expected.Error(), result.Error())
			}
		})
	}
}
