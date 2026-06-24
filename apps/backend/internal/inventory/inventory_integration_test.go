//go:build integration

package inventory_test

import (
	"context"
	"testing"

	"github.com/denden-dr/openbench/apps/backend/internal/inventory"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type InventorySuite struct {
	testutil.IntegrationSuite
	repo    inventory.InventoryRepository
	service inventory.InventoryService
}

func TestInventorySuite(t *testing.T) {
	suite.Run(t, new(InventorySuite))
}

func (s *InventorySuite) SetupTest() {
	s.IntegrationSuite.SetupTest()
	s.repo = inventory.NewRepository(s.DB)
	s.service = inventory.NewService(s.repo, s.DB)
}

func (s *InventorySuite) TestInventoryCRUD() {
	ctx := context.Background()

	// 1. Create Product
	req := &api.ProductCreate{
		Name:      "Solder Wire",
		Category:  api.ProductCreateCategoryRetail,
		Stock:     10,
		Price:     15000,
		CostPrice: 10000,
		MinStock:  2,
	}

	p, err := s.service.CreateProduct(ctx, req)
	s.Require().NoError(err)
	s.Assert().NotEmpty(p.ID)
	s.Assert().Equal("Solder Wire", p.Name)

	// 2. Get Product
	fetched, err := s.service.GetProduct(ctx, p.ID)
	s.Require().NoError(err)
	s.Assert().Equal(p.Name, fetched.Name)

	// 3. List Inventory
	list, err := s.service.ListInventory(ctx)
	s.Require().NoError(err)
	s.Assert().Len(list, 1)
	s.Assert().Equal(p.Name, list[0].Name)

	// 4. Update Product
	newName := "Super Solder Wire"
	newStock := 15
	upReq := &api.ProductUpdate{
		Name:  &newName,
		Stock: &newStock,
	}
	updated, err := s.service.UpdateProduct(ctx, p.ID, upReq)
	s.Require().NoError(err)
	s.Assert().Equal(newName, updated.Name)
	s.Assert().Equal(newStock, updated.Stock)

	// 5. Delete Product
	err = s.service.DeleteProduct(ctx, p.ID)
	s.Require().NoError(err)

	_, err = s.service.GetProduct(ctx, p.ID)
	s.Assert().Error(err)
}
