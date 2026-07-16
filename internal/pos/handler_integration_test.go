//go:build integration

package pos_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/denden-dr/OpenBench/internal/apierrors"
	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/inventory"
	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/pos"
	"github.com/denden-dr/OpenBench/internal/testutils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPosHandler_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup db container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	err = testutils.CleanTable(db, "pos_transaction_items")
	require.NoError(t, err)
	err = testutils.CleanTable(db, "pos_transactions")
	require.NoError(t, err)
	err = testutils.CleanTable(db, "products")
	require.NoError(t, err)

	// Setup layers
	posCmdRepo := pos.NewCommandRepository(db)
	posQueryRepo := pos.NewQueryRepository(db)
	invCmdRepo := inventory.NewCommandRepository(db)
	invQueryRepo := inventory.NewQueryRepository(db)
	txManager := database.NewTxManager(db)

	service := pos.NewService(posQueryRepo, posCmdRepo, invQueryRepo, invCmdRepo, txManager)
	handler := pos.NewHandler(service)

	// Setup Fiber App
	app := fiber.New(fiber.Config{
		ErrorHandler: apierrors.GlobalErrorHandler,
	})
	app.Post("/checkout", handler.Checkout)
	app.Get("/transactions", handler.GetTransactions)
	app.Get("/transactions/:id", handler.GetTransactionByID)

	// Seed Products
	prodA := &models.Product{
		ID:    uuid.New().String(),
		Name:  "Test Item A",
		Price: 10000,
		Stock: 10,
	}
	prodB := &models.Product{
		ID:    uuid.New().String(),
		Name:  "Test Item B",
		Price: 20000,
		Stock: 5,
	}
	err = invCmdRepo.Create(ctx, prodA)
	require.NoError(t, err)
	err = invCmdRepo.Create(ctx, prodB)
	require.NoError(t, err)

	var createdTxID string

	t.Run("Successful Checkout", func(t *testing.T) {
		body := models.CheckoutRequest{
			PaymentMethod: models.PaymentMethodCash,
			Items: []models.CheckoutItemRequest{
				{ProductID: prodA.ID, Quantity: 2}, // Total 20000
				{ProductID: prodB.ID, Quantity: 1}, // Total 20000
			},
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/checkout", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respData struct {
			Data struct {
				ID            string                      `json:"id"`
				PaymentMethod string                      `json:"payment_method"`
				TotalAmount   int64                       `json:"total_amount"`
				Items         []models.PosTransactionItem `json:"items"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)

		assert.NotEmpty(t, respData.Data.ID)
		assert.Equal(t, string(models.PaymentMethodCash), respData.Data.PaymentMethod)
		assert.Equal(t, int64(40000), respData.Data.TotalAmount)
		require.Len(t, respData.Data.Items, 2)

		createdTxID = respData.Data.ID

		// Verify stock deduction in database
		fetchedA, err := invQueryRepo.FindByID(ctx, prodA.ID)
		require.NoError(t, err)
		assert.Equal(t, 8, fetchedA.Stock) // 10 - 2

		fetchedB, err := invQueryRepo.FindByID(ctx, prodB.ID)
		require.NoError(t, err)
		assert.Equal(t, 4, fetchedB.Stock) // 5 - 1
	})

	t.Run("Get Transaction By ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/transactions/%s", createdTxID), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data struct {
				ID          string                      `json:"id"`
				TotalAmount int64                       `json:"total_amount"`
				Items       []models.PosTransactionItem `json:"items"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, createdTxID, respData.Data.ID)
		assert.Equal(t, int64(40000), respData.Data.TotalAmount)
		require.Len(t, respData.Data.Items, 2)
	})

	t.Run("Get Transactions List", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/transactions?limit=10&cursor=", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respData struct {
			Data []struct {
				ID          string `json:"id"`
				TotalAmount int64  `json:"total_amount"`
			} `json:"data"`
			Meta struct {
				Limit int `json:"limit"`
			} `json:"meta"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)
		assert.Equal(t, 10, respData.Meta.Limit)
		require.Len(t, respData.Data, 1)
		assert.Equal(t, createdTxID, respData.Data[0].ID)
	})

	t.Run("Checkout Failure - Insufficient Stock (Rollback verification)", func(t *testing.T) {
		body := models.CheckoutRequest{
			PaymentMethod: models.PaymentMethodQRIS,
			Items: []models.CheckoutItemRequest{
				{ProductID: prodB.ID, Quantity: 10}, // We only have 4 left
			},
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := http.NewRequest("POST", "/checkout", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		// Verify that stock was NOT decremented (transaction rolled back)
		fetchedB, err := invQueryRepo.FindByID(ctx, prodB.ID)
		require.NoError(t, err)
		assert.Equal(t, 4, fetchedB.Stock) // Still 4
	})
}
