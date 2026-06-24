package inventory

import (
	"context"
	"errors"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, p *Product) error
	GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*Product, error)
	GetByIDForUpdate(ctx context.Context, tx *sqlx.Tx, id string) (*Product, error)
	List(ctx context.Context, tx *sqlx.Tx) ([]*Product, error)
	Update(ctx context.Context, tx *sqlx.Tx, p *Product) error
	Delete(ctx context.Context, tx *sqlx.Tx, id string) error
}

type inventoryRepository struct {
	db *database.Database
}

func NewRepository(db *database.Database) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(ctx context.Context, tx *sqlx.Tx, p *Product) error {
	query := `
		INSERT INTO products (
			id, name, category, stock, price, cost_price, min_stock
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, p.ID, p.Name, p.Category, p.Stock, p.Price, p.CostPrice, p.MinStock)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, p.ID, p.Name, p.Category, p.Stock, p.Price, p.CostPrice, p.MinStock)
	}
	return err
}

func (r *inventoryRepository) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*Product, error) {
	var p Product
	query := `
		SELECT id, name, category, stock, price, cost_price, min_stock, created_at, updated_at
		FROM products WHERE id = $1
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &p, query, id)
	} else {
		err = r.db.DB.GetContext(ctx, &p, query, id)
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *inventoryRepository) GetByIDForUpdate(ctx context.Context, tx *sqlx.Tx, id string) (*Product, error) {
	var p Product
	query := `
		SELECT id, name, category, stock, price, cost_price, min_stock, created_at, updated_at
		FROM products WHERE id = $1 FOR UPDATE
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &p, query, id)
	} else {
		return nil, errors.New("GetByIDForUpdate requires an active transaction")
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *inventoryRepository) List(ctx context.Context, tx *sqlx.Tx) ([]*Product, error) {
	var products []*Product
	query := `
		SELECT id, name, category, stock, price, cost_price, min_stock, created_at, updated_at
		FROM products ORDER BY name ASC
	`
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &products, query)
	} else {
		err = r.db.DB.SelectContext(ctx, &products, query)
	}
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *inventoryRepository) Update(ctx context.Context, tx *sqlx.Tx, p *Product) error {
	query := `
		UPDATE products SET
			name = $1, category = $2, stock = $3, price = $4, cost_price = $5, min_stock = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, p.Name, p.Category, p.Stock, p.Price, p.CostPrice, p.MinStock, p.ID)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, p.Name, p.Category, p.Stock, p.Price, p.CostPrice, p.MinStock, p.ID)
	}
	return err
}

func (r *inventoryRepository) Delete(ctx context.Context, tx *sqlx.Tx, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, id)
	}
	return err
}
