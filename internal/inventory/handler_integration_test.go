//go:build integration

package inventory_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/denden-dr/OpenBench/internal/apierrors"
	"github.com/denden-dr/OpenBench/internal/inventory"
	"github.com/denden-dr/OpenBench/internal/testutils"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryHandler_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup db container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	err = testutils.CleanTable(db, "products")
	require.NoError(t, err)

	// Setup layers
	cmdRepo := inventory.NewCommandRepository(db)
	queryRepo := inventory.NewQueryRepository(db)
	service := inventory.NewService(queryRepo, cmdRepo)
	handler := inventory.NewHandler(service)

	// Setup Fiber App
	app := fiber.New(fiber.Config{
		ErrorHandler: apierrors.GlobalErrorHandler,
	})
	app.Post("/products", handler.CreateProduct)
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProductByID)
	app.Put("/products/:id", handler.UpdateProduct)
	app.Patch("/products/:id/stock", handler.AdjustStock)
	app.Delete("/products/:id", handler.DeleteProduct)

	var createdProductID string

	t.Run("Create Product", func(t *testing.T) {
		body := map[string]interface{}{
			"name":  "iPhone 15 Pro",
			"price": int64(18000000),
			"stock": 10,
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respData struct {
			Data struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Price int64  `json:"price"`
				Stock int    `json:"stock"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.NotEmpty(t, respData.Data.ID)
		assert.Equal(t, "iPhone 15 Pro", respData.Data.Name)
		assert.Equal(t, int64(18000000), respData.Data.Price)
		assert.Equal(t, 10, respData.Data.Stock)

		createdProductID = respData.Data.ID
	})

	t.Run("Get Product By ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/products/%s", createdProductID), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, createdProductID, respData.Data.ID)
		assert.Equal(t, "iPhone 15 Pro", respData.Data.Name)
	})

	t.Run("Update Product", func(t *testing.T) {
		body := map[string]interface{}{
			"name":  "iPhone 15 Pro Max",
			"price": int64(20000000),
			"stock": 8,
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PUT", fmt.Sprintf("/products/%s", createdProductID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data struct {
				Name  string `json:"name"`
				Price int64  `json:"price"`
				Stock int    `json:"stock"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, "iPhone 15 Pro Max", respData.Data.Name)
		assert.Equal(t, int64(20000000), respData.Data.Price)
		assert.Equal(t, 8, respData.Data.Stock)
	})

	t.Run("Adjust Stock", func(t *testing.T) {
		body := map[string]interface{}{
			"quantity_change": 5, // Increase stock by 5
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("PATCH", fmt.Sprintf("/products/%s/stock", createdProductID), bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data struct {
				Stock int `json:"stock"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, 13, respData.Data.Stock) // 8 + 5
	})

	t.Run("Get Products List", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/products?search=iPhone&limit=10&cursor=", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"data"`
			Meta struct {
				Limit int `json:"limit"`
			} `json:"meta"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, 10, respData.Meta.Limit)
		require.Len(t, respData.Data, 1)
		assert.Equal(t, "iPhone 15 Pro Max", respData.Data[0].Name)
	})

	t.Run("Delete Product", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/products/%s", createdProductID), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Get again, should be 404 Not Found (soft-deleted)
		getReq, err := http.NewRequest("GET", fmt.Sprintf("/products/%s", createdProductID), nil)
		require.NoError(t, err)

		getResp, err := app.Test(getReq)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})
}
