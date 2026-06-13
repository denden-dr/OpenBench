package response

import "github.com/gofiber/fiber/v2"

// APIResponse represents the standardized JSON structure for all API endpoints.
type APIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JSON sends a standardized success JSON response.
func JSON[T any](c *fiber.Ctx, statusCode int, message string, data T) error {
	return c.Status(statusCode).JSON(APIResponse[T]{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// Error sends a standardized error JSON response.
func Error(c *fiber.Ctx, statusCode int, message string, err error) error {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	return c.Status(statusCode).JSON(APIResponse[any]{
		Code:    statusCode,
		Message: message,
		Error:   errStr,
	})
}
