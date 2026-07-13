//go:build integration

package inventory_test

import (
	"context"
	"testing"

	"github.com/denden-dr/OpenBench/apps/backend/internal/inventory"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryRepository_Integration(t *testing.T) {
	ctx := context.Background()

	// Spin up PostgreSQL test container
	db, teardown, err := testutils.SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer teardown()

	cmdRepo := inventory.NewCommandRepository(db)
	queryRepo := inventory.NewQueryRepository(db)

	t.Run("CRUD Operations", func(t *testing.T) {
		err := testutils.CleanTable(db, "products")
		require.NoError(t, err)

		product := &models.Product{
			ID:    uuid.New().String(),
			Name:  "Test Product",
			Price: 150000,
			Stock: 10,
		}

		// 1. Create
		err = cmdRepo.Create(ctx, product)
		require.NoError(t, err)
		assert.NotEmpty(t, product.CreatedAt)
		assert.NotEmpty(t, product.UpdatedAt)

		// 2. FindByID
		found, err := queryRepo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, product.ID, found.ID)
		assert.Equal(t, product.Name, found.Name)
		assert.Equal(t, product.Price, found.Price)
		assert.Equal(t, product.Stock, found.Stock)

		// 3. Update
		product.Name = "Updated Test Product"
		product.Price = 200000
		product.Stock = 15
		err = cmdRepo.Update(ctx, product)
		require.NoError(t, err)

		found, err = queryRepo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, "Updated Test Product", found.Name)
		assert.Equal(t, int64(200000), found.Price)
		assert.Equal(t, 15, found.Stock)

		// 4. UpdateStock (increment/decrement)
		err = cmdRepo.UpdateStock(ctx, product.ID, -5)
		require.NoError(t, err)

		found, err = queryRepo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		assert.Equal(t, 10, found.Stock)

		// 5. Delete (soft delete)
		err = cmdRepo.Delete(ctx, product.ID)
		require.NoError(t, err)

		found, err = queryRepo.FindByID(ctx, product.ID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("FindAll Search and Pagination", func(t *testing.T) {
		err := testutils.CleanTable(db, "products")
		require.NoError(t, err)

		productsToSeed := []*models.Product{
			{ID: uuid.New().String(), Name: "Apple iPhone 15", Price: 15000000, Stock: 5},
			{ID: uuid.New().String(), Name: "Samsung Galaxy S24", Price: 14000000, Stock: 8},
			{ID: uuid.New().String(), Name: "Google Pixel 8", Price: 12000000, Stock: 12},
		}

		for _, p := range productsToSeed {
			err = cmdRepo.Create(ctx, p)
			require.NoError(t, err)
		}

		// Test FindAll with no search (returns all, sorted by name)
		list, total, err := queryRepo.FindAll(ctx, "", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, total)
		require.Len(t, list, 3)
		assert.Equal(t, "Apple iPhone 15", list[0].Name)
		assert.Equal(t, "Google Pixel 8", list[1].Name)
		assert.Equal(t, "Samsung Galaxy S24", list[2].Name)

		// Test FindAll with search
		list, total, err = queryRepo.FindAll(ctx, "galaxy", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		require.Len(t, list, 1)
		assert.Equal(t, "Samsung Galaxy S24", list[0].Name)

		// Test FindAll pagination (limit, offset)
		list, total, err = queryRepo.FindAll(ctx, "", 2, 1)
		require.NoError(t, err)
		assert.Equal(t, 3, total)
		require.Len(t, list, 2)
		assert.Equal(t, "Google Pixel 8", list[0].Name)
		assert.Equal(t, "Samsung Galaxy S24", list[1].Name)
	})
}
