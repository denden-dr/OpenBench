package middleware

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/denden-dr/openbench/apps/backend/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a global Fiber error handler that converts service, apperror,
// and validator errors into standardized HTTP JSON responses.
func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *service.AppError
	if errors.As(err, &appErr) {
		if appErr.Code >= 500 {
			if appErr.Err != nil {
				slog.Error("Internal service error occurred", "code", appErr.Code, "error", appErr.Err)
			} else {
				slog.Error("Internal service error occurred", "code", appErr.Code, "error", appErr)
			}
		}
		return c.Status(appErr.Code).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := make(map[string]string)
		for _, f := range ve {
			errs[f.Field()] = fmt.Sprintf("Field validation for '%s' failed", f.Field())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed",
			"details": errs,
		})
	}

	// Handle fiber's own errors (e.g. Bad Request on BodyParser issues or NotFound on missing routes)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		if fiberErr.Code >= 500 {
			slog.Error("Fiber server error occurred", "code", fiberErr.Code, "error", err)
		}
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"success": false,
			"error":   fiberErr.Message,
		})
	}

	slog.Error("Unhandled system error occurred", "error", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"error":   "Internal server error",
	})
}
