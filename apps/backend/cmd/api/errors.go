package main

import (
	"errors"
	"log/slog"

	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/gofiber/fiber/v3"
)

// globalErrorHandler is the custom centralized error handler for the Fiber application.
// It formats all errors, including framework-level errors (like 404, 405), to match the RFC 7807 (Problem Details) standard.
// It also masks internal details for 500 errors to prevent data leakage in production.
func globalErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	title := "Internal Server Error"
	problemType := "/errors/internal-server-error"
	detail := "An unexpected error occurred"

	switch code {
	case fiber.StatusNotFound:
		title = "Not Found"
		problemType = "/errors/not-found"
		detail = "The requested resource was not found on this server."
	case fiber.StatusMethodNotAllowed:
		title = "Method Not Allowed"
		problemType = "/errors/method-not-allowed"
		detail = "The method is not allowed for the requested URL."
	case fiber.StatusTooManyRequests:
		title = "Too Many Requests"
		problemType = "/errors/too-many-requests"
		detail = "Too many requests, please try again later."
	case fiber.StatusBadRequest:
		title = "Bad Request"
		problemType = "/errors/bad-request"
		detail = err.Error()
	case fiber.StatusUnprocessableEntity:
		title = "Unprocessable Entity"
		problemType = "/errors/unprocessable-entity"
		detail = err.Error()
	default:
		slog.ErrorContext(c.Context(), "Unhandled error caught by global ErrorHandler",
			slog.Any("error", err),
			slog.String("path", c.Path()),
			slog.String("method", c.Method()),
		)
	}

	return utils.SendProblem(c, code, problemType, title, detail)
}
