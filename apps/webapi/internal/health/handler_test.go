package health

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheckPublic(t *testing.T) {
	app := fiber.New()
	handler := NewHealthHandler(nil) // DB can be nil for public health check

	app.Get("/health", handler.HealthCheckPublic)

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
