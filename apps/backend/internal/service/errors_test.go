package service

import (
	"errors"
	"testing"

	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestMapRepositoryError(t *testing.T) {
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
			name:     "repository.ErrNotFound",
			input:    repository.ErrNotFound,
			expected: ErrTicketNotFound,
		},
		{
			name:     "repository.ErrDuplicate",
			input:    repository.ErrDuplicate,
			expected: ErrDuplicate,
		},
		{
			name:     "repository.ErrDatabaseUnavailable",
			input:    repository.ErrDatabaseUnavailable,
			expected: ErrDatabaseUnavailable,
		},
		{
			name:     "generic unknown error",
			input:    errors.New("unknown error"),
			expected: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapRepositoryError(tt.input)
			if tt.expected == nil {
				assert.NoError(t, result)
			} else {
				var appErr *AppError
				if errors.As(tt.expected, &appErr) {
					assert.ErrorIs(t, result, tt.expected)
				} else {
					assert.Equal(t, tt.expected, result)
				}
			}
		})
	}
}

func TestMapModelError(t *testing.T) {
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
			name:     "model.ErrNegativePrice",
			input:    model.ErrNegativePrice,
			expected: ErrNegativePrice,
		},
		{
			name:     "model.ErrNegativeWarranty",
			input:    model.ErrNegativeWarranty,
			expected: ErrNegativeWarranty,
		},
		{
			name:     "model.ErrPickedUpRequiresPaid",
			input:    model.ErrPickedUpRequiresPaid,
			expected: ErrInvalidPaymentStatus,
		},
		{
			name:     "model.ErrPickedUpRequiresExitDate",
			input:    model.ErrPickedUpRequiresExitDate,
			expected: ErrPickedUpMissingExitDate,
		},
		{
			name:     "model.ErrNonPickedUpCannotHaveExitDate",
			input:    model.ErrNonPickedUpCannotHaveExitDate,
			expected: ErrNonPickedUpWithDates,
		},
		{
			name:     "generic error pass-through",
			input:    errors.New("random rule"),
			expected: errors.New("random rule"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapModelError(tt.input)
			if tt.expected == nil {
				assert.NoError(t, result)
			} else {
				var appErr *AppError
				if errors.As(tt.expected, &appErr) {
					assert.ErrorIs(t, result, tt.expected)
				} else {
					assert.Equal(t, tt.expected, result)
				}
			}
		})
	}
}
