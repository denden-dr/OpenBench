package logger

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(NewMiddleware())

	var capturedCtx context.Context
	app.Get("/test", func(c fiber.Ctx) error {
		capturedCtx = c.Context()
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check response header
	reqID := resp.Header.Get(fiber.HeaderXRequestID)
	assert.NotEmpty(t, reqID)

	// Check context value
	val := capturedCtx.Value(RequestIDKey)
	assert.Equal(t, reqID, val)
}

func TestNewMiddleware_ExistingRequestID(t *testing.T) {
	app := fiber.New()
	app.Use(NewMiddleware())

	var capturedCtx context.Context
	app.Get("/test", func(c fiber.Ctx) error {
		capturedCtx = c.Context()
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(fiber.HeaderXRequestID, "custom-req-id")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check response header
	reqID := resp.Header.Get(fiber.HeaderXRequestID)
	assert.Equal(t, "custom-req-id", reqID)

	// Check context value
	val := capturedCtx.Value(RequestIDKey)
	assert.Equal(t, "custom-req-id", val)
}
