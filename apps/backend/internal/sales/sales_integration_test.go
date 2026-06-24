//go:build integration

package sales_test

import (
	"context"
	"testing"

	"github.com/denden-dr/openbench/apps/backend/internal/inventory"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/testutil"
	"github.com/denden-dr/openbench/apps/backend/internal/sales"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type SalesSuite struct {
	testutil.IntegrationSuite
	inventoryRepo inventory.InventoryRepository
	salesRepo     sales.SalesRepository
	salesService  sales.SalesService
}

func TestSalesSuite(t *testing.T) {
	suite.Run(t, new(SalesSuite))
}

func (s *SalesSuite) SetupTest() {
	s.IntegrationSuite.SetupTest()
	s.inventoryRepo = inventory.NewRepository(s.DB)
	s.salesRepo = sales.NewRepository(s.DB)
	s.salesService = sales.NewService(s.salesRepo, s.inventoryRepo, s.DB)
}

func (s *SalesSuite) TestCreateSaleSuccess() {
	ctx := context.Background()

	// 1. Create a product in inventory first
	invService := inventory.NewService(s.inventoryRepo, s.DB)
	prod, err := invService.CreateProduct(ctx, &api.ProductCreate{
		Name:      "Baterai iPhone 11",
		Category:  api.ProductCreateCategoryRetail,
		Stock:     5,
		Price:     350000,
		CostPrice: 200000,
		MinStock:  1,
	})
	s.Require().NoError(err)

	// 2. Perform checkout (Sale creation)
	req := api.CreateSaleJSONRequestBody{
		Discount:      10000,
		PaymentMethod: api.SaleCreatePaymentMethodCash,
		Items: []api.SaleCreateItem{
			{
				ProductId: uuid.MustParse(prod.ID),
				Qty:       2,
			},
		},
	}

	sale, err := s.salesService.CreateSale(ctx, &req)
	s.Require().NoError(err)
	s.Assert().NotEmpty(sale.ID)
	s.Assert().Contains(sale.InvoiceNumber, "INV-")
	s.Assert().True(sale.Total.Equal(decimal.NewFromFloat(690000.0)))
	s.Assert().Len(sale.Items, 1)

	// 3. Verify stock was deducted
	updatedProd, err := invService.GetProduct(ctx, prod.ID)
	s.Require().NoError(err)
	s.Assert().Equal(3, updatedProd.Stock) // 5 - 2 = 3

	// 4. Verify sales list
	salesList, err := s.salesService.ListSales(ctx)
	s.Require().NoError(err)
	s.Assert().Len(salesList, 1)
	s.Assert().Equal(sale.InvoiceNumber, salesList[0].InvoiceNumber)
}

func (s *SalesSuite) TestCreateSaleInsufficientStock() {
	ctx := context.Background()

	invService := inventory.NewService(s.inventoryRepo, s.DB)
	prod, err := invService.CreateProduct(ctx, &api.ProductCreate{
		Name:      "LCD iPhone 11",
		Category:  api.ProductCreateCategorySparePart,
		Stock:     1,
		Price:     600000,
		CostPrice: 400000,
		MinStock:  1,
	})
	s.Require().NoError(err)

	req := api.CreateSaleJSONRequestBody{
		Discount:      0,
		PaymentMethod: api.SaleCreatePaymentMethodQris,
		Items: []api.SaleCreateItem{
			{
				ProductId: uuid.MustParse(prod.ID),
				Qty:       2, // Requesting 2 when only 1 is in stock
			},
		},
	}

	_, err = s.salesService.CreateSale(ctx, &req)
	s.Assert().ErrorIs(err, sales.ErrInsufficientStock)

	// Verify stock remains intact
	updatedProd, err := invService.GetProduct(ctx, prod.ID)
	s.Require().NoError(err)
	s.Assert().Equal(1, updatedProd.Stock)
}
