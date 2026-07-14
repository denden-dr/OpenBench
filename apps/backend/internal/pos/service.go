package pos

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/inventory"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrInsufficientStock   = errors.New("insufficient stock for product")
)

type Service interface {
	Checkout(ctx context.Context, req models.CheckoutRequest) (*models.PosTransaction, error)
	GetTransactionByID(ctx context.Context, id string) (*models.PosTransaction, error)
	GetTransactions(ctx context.Context, limit int, cursor string) ([]models.PosTransaction, string, error)
}

type service struct {
	posQueryRepo         QueryRepository
	posCommandRepo       CommandRepository
	inventoryQueryRepo   inventory.QueryRepository
	inventoryCommandRepo inventory.CommandRepository
	txManager            database.TxManager
}

func NewService(
	posQueryRepo QueryRepository,
	posCommandRepo CommandRepository,
	inventoryQueryRepo inventory.QueryRepository,
	inventoryCommandRepo inventory.CommandRepository,
	txManager database.TxManager,
) Service {
	return &service{
		posQueryRepo:         posQueryRepo,
		posCommandRepo:       posCommandRepo,
		inventoryQueryRepo:   inventoryQueryRepo,
		inventoryCommandRepo: inventoryCommandRepo,
		txManager:            txManager,
	}
}

func (s *service) Checkout(ctx context.Context, req models.CheckoutRequest) (*models.PosTransaction, error) {
	if req.PaymentMethod != models.PaymentMethodCash && req.PaymentMethod != models.PaymentMethodQRIS {
		return nil, fmt.Errorf("%w: invalid payment method", ErrInvalidInput)
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("%w: cart cannot be empty", ErrInvalidInput)
	}

	seen := make(map[string]bool)
	for _, item := range req.Items {
		if item.ProductID == "" {
			return nil, fmt.Errorf("%w: product_id is required", ErrInvalidInput)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("%w: quantity must be greater than 0", ErrInvalidInput)
		}
		if seen[item.ProductID] {
			return nil, fmt.Errorf("%w: duplicate product in cart: %s", ErrInvalidInput, item.ProductID)
		}
		seen[item.ProductID] = true
	}

	txID := uuid.New().String()
	var totalAmount int64
	var txItems []models.PosTransactionItem

	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		totalAmount = 0
		txItems = nil

		for _, itemReq := range req.Items {
			p, err := s.inventoryQueryRepo.FindByID(txCtx, itemReq.ProductID)
			if err != nil {
				return err
			}
			if p == nil {
				return fmt.Errorf("%w: product %s not found", ErrInvalidInput, itemReq.ProductID)
			}

			if p.Stock < itemReq.Quantity {
				return fmt.Errorf("%w: %s (available: %d, requested: %d)", ErrInsufficientStock, p.Name, p.Stock, itemReq.Quantity)
			}

			err = s.inventoryCommandRepo.UpdateStock(txCtx, p.ID, -itemReq.Quantity)
			if err != nil {
				return err
			}

			itemTotal := p.Price * int64(itemReq.Quantity)
			totalAmount += itemTotal

			txItems = append(txItems, models.PosTransactionItem{
				ID:            uuid.New().String(),
				TransactionID: txID,
				ProductID:     p.ID,
				Quantity:      itemReq.Quantity,
				Price:         p.Price,
				ProductName:   p.Name,
			})
		}

		txRecord := &models.PosTransaction{
			ID:            txID,
			PaymentMethod: req.PaymentMethod,
			TotalAmount:   totalAmount,
			Items:         txItems,
		}

		return s.posCommandRepo.Create(txCtx, txRecord)
	})

	if err != nil {
		return nil, err
	}

	tx, err := s.GetTransactionByID(ctx, txID)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "POS checkout completed",
		slog.String("transaction_id", tx.ID),
		slog.String("payment_method", string(tx.PaymentMethod)),
		slog.Int64("total_amount", tx.TotalAmount),
		slog.Int("item_count", len(tx.Items)),
	)

	return tx, nil
}

func (s *service) GetTransactionByID(ctx context.Context, id string) (*models.PosTransaction, error) {
	t, err := s.posQueryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrTransactionNotFound
	}
	return t, nil
}

func (s *service) GetTransactions(ctx context.Context, limit int, cursor string) ([]models.PosTransaction, string, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > utils.MaxLimit {
		limit = utils.MaxLimit
	}
	return s.posQueryRepo.FindAll(ctx, limit, cursor)
}
