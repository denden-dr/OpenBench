package inventory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
)

type InventoryService interface {
	CreateProduct(ctx context.Context, req *api.ProductCreate) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	ListInventory(ctx context.Context) ([]*Product, error)
	UpdateProduct(ctx context.Context, id string, req *api.ProductUpdate) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type inventoryService struct {
	repo InventoryRepository
	db   *database.Database
}

func NewService(repo InventoryRepository, db *database.Database) InventoryService {
	return &inventoryService{
		repo: repo,
		db:   db,
	}
}

func (s *inventoryService) CreateProduct(ctx context.Context, req *api.ProductCreate) (*Product, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}
	if req.Price < 0 || req.CostPrice < 0 {
		return nil, fmt.Errorf("%w: price and cost price must be non-negative", ErrInvalidInput)
	}

	p := &Product{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Category:  string(req.Category),
		Stock:     req.Stock,
		Price:     decimal.NewFromFloat32(req.Price),
		CostPrice: decimal.NewFromFloat32(req.CostPrice),
		MinStock:  req.MinStock,
	}

	if err := s.repo.Create(ctx, nil, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *inventoryService) GetProduct(ctx context.Context, id string) (*Product, error) {
	p, err := s.repo.GetByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *inventoryService) ListInventory(ctx context.Context) ([]*Product, error) {
	return s.repo.List(ctx, nil)
}

func (s *inventoryService) UpdateProduct(ctx context.Context, id string, req *api.ProductUpdate) (*Product, error) {
	p, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Category != nil {
		p.Category = string(*req.Category)
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	if req.Price != nil {
		if *req.Price < 0 {
			return nil, fmt.Errorf("%w: price must be non-negative", ErrInvalidInput)
		}
		p.Price = decimal.NewFromFloat32(*req.Price)
	}
	if req.CostPrice != nil {
		if *req.CostPrice < 0 {
			return nil, fmt.Errorf("%w: cost price must be non-negative", ErrInvalidInput)
		}
		p.CostPrice = decimal.NewFromFloat32(*req.CostPrice)
	}
	if req.MinStock != nil {
		p.MinStock = *req.MinStock
	}

	if err := s.repo.Update(ctx, nil, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *inventoryService) DeleteProduct(ctx context.Context, id string) error {
	_, err := s.GetProduct(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, nil, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrProductReferenced
		}
		return err
	}
	return nil
}
