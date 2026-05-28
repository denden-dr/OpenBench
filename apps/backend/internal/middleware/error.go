package middleware

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/denden-dr/openbench/apps/backend/internal/dto"
	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a global Fiber error handler that converts service, apperror,
// and validator errors into RFC 7807 Problem Details HTTP JSON responses.
func ErrorHandler(c *fiber.Ctx, err error) error {
	c.Set(fiber.HeaderContentType, "application/problem+json")

	// 1. service.AppError
	var appErr *service.AppError
	if errors.As(err, &appErr) {
		status := appErr.Code
		if status == 0 {
			status = 500
		}

		if status >= 500 {
			if appErr.Err != nil {
				slog.Error("Internal service error occurred", "code", status, "error", appErr.Err)
			} else {
				slog.Error("Internal service error occurred", "code", status, "error", appErr)
			}
		}

		var typeURI string
		var title string
		switch status {
		case fiber.StatusBadRequest:
			typeURI = "https://openbench.denden.com/errors/validation-failed"
			title = "Validation Failed"
		case fiber.StatusNotFound:
			typeURI = "https://openbench.denden.com/errors/not-found"
			title = "Not Found"
		case fiber.StatusConflict:
			typeURI = "https://openbench.denden.com/errors/conflict"
			title = "Conflict"
		case fiber.StatusServiceUnavailable:
			typeURI = "https://openbench.denden.com/errors/database-unavailable"
			title = "Database Unavailable"
		default:
			typeURI = "https://openbench.denden.com/errors/internal-error"
			title = "Internal Server Error"
		}

		detail := appErr.Message

		return c.Status(status).JSON(dto.ProblemDetails{
			Type:     typeURI,
			Title:    title,
			Status:   status,
			Detail:   detail,
			Instance: c.Path(),
		})
	}

	// 2. validator.ValidationErrors
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		invalidParams := make(map[string]string)
		for _, f := range ve {
			invalidParams[f.Field()] = fmt.Sprintf("Field validation for '%s' failed", f.Field())
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.ProblemDetails{
			Type:          "https://openbench.denden.com/errors/validation-failed",
			Title:         "Validation Failed",
			Status:        fiber.StatusBadRequest,
			Detail:        "Validation failed for one or more fields",
			Instance:      c.Path(),
			InvalidParams: invalidParams,
		})
	}

	// 3. fiber.Error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		status := fiberErr.Code
		if status >= 500 {
			slog.Error("Fiber server error occurred", "code", status, "error", err)
		}

		var typeURI string
		var title string
		switch status {
		case fiber.StatusBadRequest:
			typeURI = "https://openbench.denden.com/errors/validation-failed"
			title = "Validation Failed"
		case fiber.StatusNotFound:
			typeURI = "https://openbench.denden.com/errors/not-found"
			title = "Not Found"
		default:
			typeURI = "https://openbench.denden.com/errors/internal-error"
			title = "Internal Server Error"
		}

		return c.Status(status).JSON(dto.ProblemDetails{
			Type:     typeURI,
			Title:    title,
			Status:   status,
			Detail:   fiberErr.Message,
			Instance: c.Path(),
		})
	}

	// 4. Default unhandled error
	slog.Error("Unhandled system error occurred", "error", err)
	return c.Status(fiber.StatusInternalServerError).JSON(dto.ProblemDetails{
		Type:     "https://openbench.denden.com/errors/internal-error",
		Title:    "Internal Server Error",
		Status:   fiber.StatusInternalServerError,
		Detail:   "Internal server error",
		Instance: c.Path(),
	})
}
