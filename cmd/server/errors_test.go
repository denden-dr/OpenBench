package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/denden-dr/OpenBench/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalErrorHandler(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: globalErrorHandler,
	})

	// Test route for throwing a standard fiber error
	app.Get("/fiber-error", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "custom not found msg")
	})

	// Test route for throwing a random/generic error
	app.Get("/generic-error", func(c fiber.Ctx) error {
		return errors.New("this is a database crash or secret key error")
	})

	// Test route for throwing a bad request error
	app.Get("/bad-request", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "invalid query parameter 'id'")
	})

	// 1. Verify standard fiber 404 Not Found error
	req, err := http.NewRequest("GET", "/fiber-error", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "application/problem+json", resp.Header.Get("Content-Type"))

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var problem utils.ProblemDetail
	err = json.Unmarshal(bodyBytes, &problem)
	require.NoError(t, err)

	assert.Contains(t, problem.Type, "/errors/not-found")
	assert.Equal(t, "Not Found", problem.Title)
	assert.Equal(t, fiber.StatusNotFound, problem.Status)
	assert.Equal(t, "The requested resource was not found on this server.", problem.Detail)
	assert.Equal(t, "/fiber-error", problem.Instance)

	// 2. Verify 500 error is masked
	req, err = http.NewRequest("GET", "/generic-error", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	bodyBytes, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(bodyBytes, &problem)
	require.NoError(t, err)

	assert.Contains(t, problem.Type, "/errors/internal-server-error")
	assert.Equal(t, "Internal Server Error", problem.Title)
	assert.Equal(t, "An unexpected error occurred", problem.Detail) // Masked!

	// 3. Verify 400 Bad Request error detail is exposed (validation message)
	req, err = http.NewRequest("GET", "/bad-request", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	bodyBytes, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(bodyBytes, &problem)
	require.NoError(t, err)

	assert.Contains(t, problem.Type, "/errors/bad-request")
	assert.Equal(t, "Bad Request", problem.Title)
	assert.Equal(t, "invalid query parameter 'id'", problem.Detail) // Exposed!
}

func TestRecoverMiddleware(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: globalErrorHandler,
	})

	app.Use(recover.New())

	// Test route that panics
	app.Get("/panic", func(c fiber.Ctx) error {
		panic("boom!")
	})

	req, err := http.NewRequest("GET", "/panic", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Recover should catch panic and pass it to ErrorHandler, resulting in 500
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var problem utils.ProblemDetail
	err = json.Unmarshal(bodyBytes, &problem)
	require.NoError(t, err)

	assert.Contains(t, problem.Type, "/errors/internal-server-error")
	assert.Equal(t, "Internal Server Error", problem.Title)
	assert.Equal(t, "An unexpected error occurred", problem.Detail) // Masked!
}
