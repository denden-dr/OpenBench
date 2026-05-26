package model

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func NewModelError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

var (
	ErrNegativePrice                 = NewModelError(1, "price cannot be negative")
	ErrNegativeWarranty              = NewModelError(2, "warranty days cannot be negative")
	ErrPickedUpRequiresPaid          = NewModelError(3, "picked up ticket must be paid")
	ErrPickedUpRequiresExitDate      = NewModelError(4, "picked up ticket must have an exit date")
	ErrNonPickedUpCannotHaveExitDate = NewModelError(5, "non-picked up ticket cannot have exit date")
	ErrInvalidStatusTransition       = NewModelError(6, "invalid ticket status transition")
)
