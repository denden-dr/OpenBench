package inventory

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidInput    = errors.New("invalid input data")
)

type CreateProductRequest struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

type UpdateProductRequest struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

type Service interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error)
	UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*models.Product, error)
	AdjustStock(ctx context.Context, id string, quantityChange int) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProducts(ctx context.Context, search string, limit, offset int) ([]models.Product, int, error)
	DeleteProduct(ctx context.Context, id string) error
}

type service struct {
	queryRepo   QueryRepository
	commandRepo CommandRepository
}

func NewService(queryRepo QueryRepository, commandRepo CommandRepository) Service {
	return &service{
		queryRepo:   queryRepo,
		commandRepo: commandRepo,
	}
}

func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}
	if req.Stock < 0 {
		return nil, fmt.Errorf("%w: stock cannot be negative", ErrInvalidInput)
	}

	p := &models.Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: req.Price,
		Stock: req.Stock,
	}

	if err := s.commandRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*models.Product, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}
	if req.Stock < 0 {
		return nil, fmt.Errorf("%w: stock cannot be negative", ErrInvalidInput)
	}

	p, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}

	p.Name = name
	p.Price = req.Price
	p.Stock = req.Stock

	if err := s.commandRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) AdjustStock(ctx context.Context, id string, quantityChange int) error {
	p, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrProductNotFound
	}

	// Prevent negative stock
	if p.Stock+quantityChange < 0 {
		return fmt.Errorf("%w: stock cannot be less than 0", ErrInvalidInput)
	}

	return s.commandRepo.UpdateStock(ctx, id, quantityChange)
}

func (s *service) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	p, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	return p, nil
}

func (s *service) GetProducts(ctx context.Context, search string, limit, offset int) ([]models.Product, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.queryRepo.FindAll(ctx, search, limit, offset)
}

func (s *service) DeleteProduct(ctx context.Context, id string) error {
	p, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrProductNotFound
	}
	return s.commandRepo.Delete(ctx, id)
}
