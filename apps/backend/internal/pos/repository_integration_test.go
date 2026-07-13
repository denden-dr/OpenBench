//go:build integration

package pos_test

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/inventory"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/pos"
	"github.com/denden-dr/OpenBench/apps/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPosRepository_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup db container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	cmdRepo := pos.NewCommandRepository(db)
	queryRepo := pos.NewQueryRepository(db)
	invCmdRepo := inventory.NewCommandRepository(db)

	t.Run("Create and Retrieve Transactions", func(t *testing.T) {
		err := testutils.CleanTable(db, "pos_transaction_items")
		require.NoError(t, err)
		err = testutils.CleanTable(db, "pos_transactions")
		require.NoError(t, err)
		err = testutils.CleanTable(db, "products")
		require.NoError(t, err)

		// Seed a product first for foreign key constraint
		prod := &models.Product{
			ID:    uuid.New().String(),
			Name:  "POS Keyboard",
			Price: 450000,
			Stock: 20,
		}
		err = invCmdRepo.Create(ctx, prod)
		require.NoError(t, err)

		tx := &models.PosTransaction{
			ID:            uuid.New().String(),
			PaymentMethod: models.PaymentMethodCash,
			TotalAmount:   450000,
			Items: []models.PosTransactionItem{
				{
					ID:        uuid.New().String(),
					ProductID: prod.ID,
					Quantity:  1,
					Price:     450000,
				},
			},
		}

		// 1. Create Transaction
		err = cmdRepo.Create(ctx, tx)
		require.NoError(t, err)
		assert.NotEmpty(t, tx.CreatedAt)

		// 2. FindByID
		found, err := queryRepo.FindByID(ctx, tx.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, tx.ID, found.ID)
		assert.Equal(t, models.PaymentMethodCash, found.PaymentMethod)
		assert.Equal(t, int64(450000), found.TotalAmount)
		require.Len(t, found.Items, 1)

		item := found.Items[0]
		assert.Equal(t, tx.Items[0].ID, item.ID)
		assert.Equal(t, prod.ID, item.ProductID)
		assert.Equal(t, 1, item.Quantity)
		assert.Equal(t, int64(450000), item.Price)
		assert.Equal(t, "POS Keyboard", item.ProductName) // Joined product name check
	})

	t.Run("FindAll Transactions", func(t *testing.T) {
		err := testutils.CleanTable(db, "pos_transaction_items")
		require.NoError(t, err)
		err = testutils.CleanTable(db, "pos_transactions")
		require.NoError(t, err)

		tx1 := &models.PosTransaction{
			ID:            uuid.New().String(),
			PaymentMethod: models.PaymentMethodCash,
			TotalAmount:   100000,
		}
		tx2 := &models.PosTransaction{
			ID:            uuid.New().String(),
			PaymentMethod: models.PaymentMethodQRIS,
			TotalAmount:   200000,
		}

		err = cmdRepo.Create(ctx, tx1)
		require.NoError(t, err)
		err = cmdRepo.Create(ctx, tx2)
		require.NoError(t, err)

		list, total, err := queryRepo.FindAll(ctx, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		require.Len(t, list, 2)
		// Ordered by created_at DESC
		assert.Equal(t, tx2.ID, list[0].ID)
		assert.Equal(t, tx1.ID, list[1].ID)

		// Test limit and offset
		list, total, err = queryRepo.FindAll(ctx, 1, 1)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		require.Len(t, list, 1)
		assert.Equal(t, tx1.ID, list[0].ID)
	})
}
