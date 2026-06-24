package sales

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/inventory"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrSaleNotFound      = errors.New("sale not found")
	ErrInsufficientStock = errors.New("insufficient stock for product")
	ErrInvalidInput      = errors.New("invalid input data")
)

type SalesService interface {
	CreateSale(ctx context.Context, req *api.CreateSaleJSONRequestBody) (*Sale, error)
	ListSales(ctx context.Context) ([]*Sale, error)
}

type salesService struct {
	repo          SalesRepository
	inventoryRepo inventory.InventoryRepository
	db            *database.Database
}

func NewService(repo SalesRepository, inventoryRepo inventory.InventoryRepository, db *database.Database) SalesService {
	return &salesService{
		repo:          repo,
		inventoryRepo: inventoryRepo,
		db:            db,
	}
}

func (s *salesService) CreateSale(ctx context.Context, req *api.CreateSaleJSONRequestBody) (*Sale, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("%w: sale must contain at least one item", ErrInvalidInput)
	}

	tx, err := s.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Generate Invoice Number
	invoiceNum, err := s.generateInvoiceNumber(ctx, tx)
	if err != nil {
		return nil, err
	}

	saleID := uuid.New().String()
	sale := &Sale{
		ID:            saleID,
		InvoiceNumber: invoiceNum,
		Discount:      float64(req.Discount),
		PaymentMethod: string(req.PaymentMethod),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	var saleItems []SaleItem
	var subtotal float64

	// 2. Process Sale Items and Stock Deduction
	for _, item := range req.Items {
		prodID := item.ProductId.String()
		product, err := s.inventoryRepo.GetByIDForUpdate(ctx, tx, prodID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product %s: %w", prodID, err)
		}

		if product.Stock < item.Qty {
			return nil, fmt.Errorf("%w: %s (available: %d, requested: %d)", ErrInsufficientStock, product.Name, product.Stock, item.Qty)
		}

		// Deduct stock
		product.Stock -= item.Qty
		if err := s.inventoryRepo.Update(ctx, tx, product); err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}

		itemPrice := product.Price
		subtotal += itemPrice * float64(item.Qty)

		saleItem := SaleItem{
			ID:        uuid.New().String(),
			SaleID:    saleID,
			ProductID: prodID,
			Name:      product.Name,
			Price:     itemPrice,
			Qty:       item.Qty,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		saleItems = append(saleItems, saleItem)
	}

	if sale.Discount > subtotal {
		return nil, fmt.Errorf("%w: discount cannot exceed subtotal", ErrInvalidInput)
	}

	sale.Subtotal = subtotal
	sale.Total = subtotal - sale.Discount
	sale.Items = saleItems

	// 3. Insert Sale
	if err := s.repo.Create(ctx, tx, sale); err != nil {
		return nil, err
	}

	// 4. Insert Sale Items
	for _, si := range saleItems {
		siCopy := si
		if err := s.repo.CreateItem(ctx, tx, &siCopy); err != nil {
			return nil, fmt.Errorf("failed to insert sale item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *salesService) ListSales(ctx context.Context) ([]*Sale, error) {
	return s.repo.List(ctx, nil)
}

func (s *salesService) generateInvoiceNumber(ctx context.Context, tx *sqlx.Tx) (string, error) {
	prefix := fmt.Sprintf("INV-%s-", time.Now().Format("200601"))
	nextNum, err := s.repo.GetNextInvoiceSequence(ctx, tx, prefix)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, nextNum), nil
}
