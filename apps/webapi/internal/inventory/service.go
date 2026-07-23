package inventory

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidInput    = errors.New("invalid input data")
)

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required,min=0"`
	Stock int    `json:"stock" validate:"min=0"`
}

type UpdateProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required,min=0"`
	Stock int    `json:"stock" validate:"min=0"`
}

type Service interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error)
	UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*models.Product, error)
	AdjustStock(ctx context.Context, id string, quantityChange int) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProducts(ctx context.Context, search string, limit int, cursor string) ([]models.Product, string, error)
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

func validateProductInput(name string, price int64, stock int) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return "", fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if price < 0 {
		return "", fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}
	if stock < 0 {
		return "", fmt.Errorf("%w: stock cannot be negative", ErrInvalidInput)
	}
	return trimmedName, nil
}

func (s *service) getProductOrError(ctx context.Context, id string) (*models.Product, error) {
	p, err := s.queryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	return p, nil
}

func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {
	name, err := validateProductInput(req.Name, req.Price, req.Stock)
	if err != nil {
		return nil, err
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

	slog.InfoContext(ctx, "Product created",
		slog.String("product_id", p.ID),
		slog.String("name", p.Name),
		slog.Int64("price", p.Price),
		slog.Int("stock", p.Stock),
	)

	return p, nil
}

func (s *service) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*models.Product, error) {
	name, err := validateProductInput(req.Name, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}

	p, err := s.getProductOrError(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Name = name
	p.Price = req.Price
	p.Stock = req.Stock

	if err := s.commandRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Product updated",
		slog.String("product_id", p.ID),
		slog.String("name", p.Name),
		slog.Int64("price", p.Price),
		slog.Int("stock", p.Stock),
	)

	return p, nil
}

func (s *service) AdjustStock(ctx context.Context, id string, quantityChange int) error {
	p, err := s.getProductOrError(ctx, id)
	if err != nil {
		return err
	}

	// Prevent negative stock
	if p.Stock+quantityChange < 0 {
		return fmt.Errorf("%w: stock cannot be less than 0", ErrInvalidInput)
	}

	if err := s.commandRepo.UpdateStock(ctx, id, quantityChange); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Product stock adjusted",
		slog.String("product_id", id),
		slog.String("product_name", p.Name),
		slog.Int("quantity_change", quantityChange),
		slog.Int("stock_before", p.Stock),
		slog.Int("stock_after", p.Stock+quantityChange),
	)

	return nil
}

func (s *service) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	return s.getProductOrError(ctx, id)
}

func (s *service) GetProducts(ctx context.Context, search string, limit int, cursor string) ([]models.Product, string, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > utils.MaxLimit {
		limit = utils.MaxLimit
	}
	return s.queryRepo.FindAll(ctx, search, limit, cursor)
}

func (s *service) DeleteProduct(ctx context.Context, id string) error {
	p, err := s.getProductOrError(ctx, id)
	if err != nil {
		return err
	}
	if err := s.commandRepo.Delete(ctx, id); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Product deleted",
		slog.String("product_id", id),
		slog.String("product_name", p.Name),
	)

	return nil
}
