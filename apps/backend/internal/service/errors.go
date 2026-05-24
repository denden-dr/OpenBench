package service

import (
	"errors"
	"fmt"

	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/denden-dr/openbench/apps/backend/internal/repository"
)

type AppError struct {
	Code    int
	Message string
	Err     error // Holds the underlying raw error for logging/debugging
}

func (a *AppError) Error() string {
	if a.Err != nil {
		return fmt.Sprintf("%d: %s: %v", a.Code, a.Message, a.Err)
	}
	return fmt.Sprintf("%d: %s", a.Code, a.Message)
}

func (a *AppError) Unwrap() error {
	return a.Err
}

// Is compares AppErrors based on Code and Message, ignoring the underlying Err.
// This allows errors.Is(err, ErrInternal) to succeed even if the instance has a wrapped DB error.
func (a *AppError) Is(target error) bool {
	var t *AppError
	if errors.As(target, &t) {
		return a.Code == t.Code && a.Message == t.Message
	}
	return false
}

func NewAppError(code int, message string) error {
	return &AppError{Code: code, Message: message}
}

var (
	ErrTicketNotFound          = NewAppError(404, "ticket not found")
	ErrInvalidPaymentStatus    = NewAppError(400, "a picked up ticket must be paid")
	ErrNegativePrice           = NewAppError(400, "price cannot be negative")
	ErrNegativeWarranty        = NewAppError(400, "warranty days cannot be negative")
	ErrNonPickedUpWithDates    = NewAppError(400, "non-picked up ticket cannot have exit date")
	ErrPickedUpMissingExitDate = NewAppError(400, "picked up ticket must have an exit date")

	ErrDuplicate           = NewAppError(409, "resource already exists")
	ErrDatabaseUnavailable = NewAppError(503, "database is currently unavailable")
	ErrInternal            = NewAppError(500, "internal server error")
)

// MapModelError maps model business rule errors to service AppErrors
func MapModelError(err error) error {
	if err == nil {
		return nil
	}
	var modelErr *model.Error
	if errors.As(err, &modelErr) {
		switch {
		case errors.Is(modelErr, model.ErrNegativePrice):
			return ErrNegativePrice
		case errors.Is(modelErr, model.ErrNegativeWarranty):
			return ErrNegativeWarranty
		case errors.Is(modelErr, model.ErrPickedUpRequiresPaid):
			return ErrInvalidPaymentStatus
		case errors.Is(modelErr, model.ErrPickedUpRequiresExitDate):
			return ErrPickedUpMissingExitDate
		case errors.Is(modelErr, model.ErrNonPickedUpCannotHaveExitDate):
			return ErrNonPickedUpWithDates
		}
	}
	return err
}

// MapRepositoryError maps repository-level errors to service AppErrors and wraps the raw error
func MapRepositoryError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrTicketNotFound
	case errors.Is(err, repository.ErrDuplicate):
		return ErrDuplicate
	case errors.Is(err, repository.ErrDatabaseUnavailable):
		return &AppError{Code: 503, Message: "database is currently unavailable", Err: err}
	default:
		return &AppError{Code: 500, Message: "internal server error", Err: err}
	}
}
