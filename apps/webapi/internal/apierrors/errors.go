package apierrors

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/denden-dr/OpenBench/internal/auth"
	"github.com/denden-dr/OpenBench/internal/inventory"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/pos"
	"github.com/denden-dr/OpenBench/internal/ticket"
	"github.com/denden-dr/OpenBench/internal/utils"
	"github.com/denden-dr/OpenBench/internal/warranty"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// GlobalErrorHandler is the custom centralized error handler for the Fiber application.
// It formats all errors, including framework-level errors (like 404, 405), to match the RFC 7807 (Problem Details) standard.
// It also masks internal details for 500 errors to prevent data leakage in production.
func GlobalErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	title := "Internal Server Error"
	problemType := "/errors/internal-server-error"
	detail := "An unexpected error occurred"

	var ve validator.ValidationErrors

	switch {
	case errors.As(err, &ve):
		code = fiber.StatusBadRequest
		title = "Validation Failed"
		problemType = "/errors/bad-request"
		lang := c.Get("Accept-Language")
		locale := "en"
		if strings.HasPrefix(strings.ToLower(lang), "id") {
			locale = "id"
		}
		translatedErrs := utils.TranslateValidationErrors(ve, locale)
		detail = "Validation failed: " + strings.Join(translatedErrs, ", ")

	case errors.Is(err, auth.ErrInvalidCredentials) || errors.Is(err, auth.ErrInvalidToken):
		code = fiber.StatusUnauthorized
		title = "Unauthorized"
		problemType = "/errors/unauthorized"
		detail = err.Error()

	case errors.Is(err, auth.ErrUserNotFound) ||
		errors.Is(err, inventory.ErrProductNotFound) ||
		errors.Is(err, pos.ErrTransactionNotFound) ||
		errors.Is(err, ticket.ErrTicketNotFound) ||
		errors.Is(err, warranty.ErrWarrantyNotFound) ||
		errors.Is(err, warranty.ErrClaimNotFound):
		code = fiber.StatusNotFound
		title = "Not Found"
		problemType = "/errors/not-found"
		detail = err.Error()

	case errors.Is(err, pos.ErrInsufficientStock):
		code = fiber.StatusConflict
		title = "Conflict - Insufficient Stock"
		problemType = "/errors/conflict"
		detail = err.Error()

	case errors.Is(err, inventory.ErrInvalidInput) ||
		errors.Is(err, pos.ErrInvalidInput) ||
		errors.Is(err, warranty.ErrInvalidInput) ||
		errors.Is(err, warranty.ErrWarrantyNotActive) ||
		errors.Is(err, ticket.ErrInvalidInput) ||
		errors.Is(err, models.ErrMissingCustomerName) ||
		errors.Is(err, models.ErrMissingCustomerPhone) ||
		errors.Is(err, models.ErrMissingDeviceBrand) ||
		errors.Is(err, models.ErrMissingDeviceModel) ||
		errors.Is(err, models.ErrMissingIssueDescription) ||
		errors.Is(err, models.ErrNegativeCost) ||
		errors.Is(err, models.ErrNegativeWarrantyDays):
		code = fiber.StatusBadRequest
		title = "Bad Request"
		problemType = "/errors/bad-request"
		detail = err.Error()

	default:
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
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
			}
		} else {
			type stackTracer interface {
				StackTrace() string
			}
			var st stackTracer
			if errors.As(err, &st) {
				slog.ErrorContext(c.Context(), "Unhandled error caught by global ErrorHandler",
					slog.Any("error", err),
					slog.String("stack_trace", st.StackTrace()),
					slog.String("path", c.Path()),
					slog.String("method", c.Method()),
				)
			} else {
				slog.ErrorContext(c.Context(), "Unhandled error caught by global ErrorHandler",
					slog.Any("error", err),
					slog.String("path", c.Path()),
					slog.String("method", c.Method()),
				)
			}
		}
	}

	return utils.SendProblem(c, code, problemType, title, detail)
}
