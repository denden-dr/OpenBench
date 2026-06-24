package sales

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/inventory"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
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

	// 1. Accumulate quantities by product ID and extract unique sorted IDs
	itemQtyMap := make(map[string]int)
	for _, item := range req.Items {
		prodID := item.ProductId.String()
		itemQtyMap[prodID] += item.Qty
	}

	var productIDs []string
	for prodID := range itemQtyMap {
		productIDs = append(productIDs, prodID)
	}
	sort.Strings(productIDs)

	saleID := uuid.New().String()
	sale := &Sale{
		ID:            saleID,
		Discount:      decimal.NewFromFloat32(req.Discount),
		PaymentMethod: string(req.PaymentMethod),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	var saleItems []SaleItem
	subtotal := decimal.NewFromInt(0)

	// 2. Process Sale Items and Stock Deduction in deterministic order
	for _, prodID := range productIDs {
		qty := itemQtyMap[prodID]
		product, err := s.inventoryRepo.GetByIDForUpdate(ctx, tx, prodID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%w: product %s not found", ErrInvalidInput, prodID)
			}
			return nil, fmt.Errorf("failed to fetch product %s: %w", prodID, err)
		}

		if product.Stock < qty {
			return nil, fmt.Errorf("%w: %s (available: %d, requested: %d)", ErrInsufficientStock, product.Name, product.Stock, qty)
		}

		// Deduct stock
		product.Stock -= qty
		if err := s.inventoryRepo.Update(ctx, tx, product); err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}

		itemPrice := product.Price
		subtotal = subtotal.Add(itemPrice.Mul(decimal.NewFromInt(int64(qty))))

		saleItem := SaleItem{
			ID:        uuid.New().String(),
			SaleID:    saleID,
			ProductID: prodID,
			Name:      product.Name,
			Price:     itemPrice,
			Qty:       qty,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		saleItems = append(saleItems, saleItem)
	}

	if sale.Discount.GreaterThan(subtotal) {
		return nil, fmt.Errorf("%w: discount cannot exceed subtotal", ErrInvalidInput)
	}

	sale.Subtotal = subtotal
	sale.Total = subtotal.Sub(sale.Discount)
	sale.Items = saleItems

	// 3. Generate Invoice Number (moved down to prevent bottleneck)
	invoiceNum, err := s.generateInvoiceNumber(ctx, tx)
	if err != nil {
		return nil, err
	}
	sale.InvoiceNumber = invoiceNum

	// 4. Insert Sale
	if err := s.repo.Create(ctx, tx, sale); err != nil {
		return nil, err
	}

	// 5. Insert Sale Items
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
