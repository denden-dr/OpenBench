package apierrors

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/auth"
	"github.com/denden-dr/OpenBench/apps/backend/internal/inventory"
	"github.com/denden-dr/OpenBench/apps/backend/internal/pos"
	"github.com/denden-dr/OpenBench/apps/backend/internal/ticket"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/denden-dr/OpenBench/apps/backend/internal/warranty"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalErrorHandler_DomainErrors(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: GlobalErrorHandler,
	})

	app.Get("/auth-error", func(c fiber.Ctx) error {
		return auth.ErrInvalidCredentials
	})

	app.Get("/notfound-error", func(c fiber.Ctx) error {
		return inventory.ErrProductNotFound
	})

	app.Get("/conflict-error", func(c fiber.Ctx) error {
		return pos.ErrInsufficientStock
	})

	app.Get("/badrequest-error", func(c fiber.Ctx) error {
		return ticket.ErrInvalidInput
	})

	app.Get("/warranty-notactive", func(c fiber.Ctx) error {
		return warranty.ErrWarrantyNotActive
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth-error", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/notfound-error", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Conflict", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/conflict-error", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("BadRequest", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/badrequest-error", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("WarrantyNotActive", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/warranty-notactive", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Validation Failed - English", func(t *testing.T) {
		type valStruct struct {
			Name string `validate:"required"`
		}
		app.Post("/validate-en", func(c fiber.Ctx) error {
			var req valStruct
			_ = c.Bind().JSON(&req)
			return utils.ValidateStruct(req)
		})

		req := httptest.NewRequest("POST", "/validate-en", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Language", "en")
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var problem utils.ProblemDetail
		err = json.Unmarshal(bodyBytes, &problem)
		require.NoError(t, err)
		assert.Contains(t, problem.Detail, "Validation failed: Name is a required field")
	})

	t.Run("Validation Failed - Indonesian", func(t *testing.T) {
		type valStruct struct {
			Name string `validate:"required"`
		}
		app.Post("/validate-id", func(c fiber.Ctx) error {
			var req valStruct
			_ = c.Bind().JSON(&req)
			return utils.ValidateStruct(req)
		})

		req := httptest.NewRequest("POST", "/validate-id", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Language", "id")
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var problem utils.ProblemDetail
		err = json.Unmarshal(bodyBytes, &problem)
		require.NoError(t, err)
		assert.Contains(t, problem.Detail, "Validation failed: Name wajib diisi")
	})

	t.Run("Fiber Errors", func(t *testing.T) {
		app.Get("/fiber-badrequest", func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusBadRequest, "some msg")
		})
		app.Get("/fiber-methodnotallowed", func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusMethodNotAllowed, "some msg")
		})
		app.Get("/fiber-unprocessable", func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "some msg")
		})
		app.Get("/fiber-toomany", func(c fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "some msg")
		})

		for _, tc := range []struct {
			path string
			code int
		}{
			{"/fiber-badrequest", http.StatusBadRequest},
			{"/fiber-methodnotallowed", http.StatusMethodNotAllowed},
			{"/fiber-unprocessable", http.StatusUnprocessableEntity},
			{"/fiber-toomany", http.StatusTooManyRequests},
		} {
			req := httptest.NewRequest("GET", tc.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tc.code, resp.StatusCode)
			resp.Body.Close()
		}
	})

	t.Run("Unhandled error logs and returns 500", func(t *testing.T) {
		app.Get("/unhandled", func(c fiber.Ctx) error {
			return errors.New("something went wrong internally")
		})

		req := httptest.NewRequest("GET", "/unhandled", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
