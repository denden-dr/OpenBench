package inventory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/jmoiron/sqlx"
)

type QueryRepository interface {
	FindByID(ctx context.Context, id string) (*models.Product, error)
	FindAll(ctx context.Context, search string, limit int, cursor string) ([]models.Product, string, error)
}

type CommandRepository interface {
	Create(ctx context.Context, p *models.Product) error
	Update(ctx context.Context, p *models.Product) error
	UpdateStock(ctx context.Context, id string, quantityChange int) error
	Delete(ctx context.Context, id string) error
}

type sqlQueryRepository struct {
	db *sqlx.DB
}

type sqlCommandRepository struct {
	db *sqlx.DB
}

func NewQueryRepository(db *sqlx.DB) QueryRepository {
	return &sqlQueryRepository{db: db}
}

func NewCommandRepository(db *sqlx.DB) CommandRepository {
	return &sqlCommandRepository{db: db}
}

func (r *sqlQueryRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var p models.Product
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *sqlQueryRepository) FindAll(ctx context.Context, search string, limit int, cursor string) ([]models.Product, string, error) {
	var selectQuery = `
		SELECT id, name, price, stock, created_at, updated_at
		FROM products
		WHERE deleted_at IS NULL
	`

	var conditions []string
	var args []interface{}
	argCount := 1

	if search != "" {
		searchPattern := "%" + search + "%"
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argCount))
		args = append(args, searchPattern)
		argCount++
	}

	if cursor != "" {
		cursorTime, cursorID, err := utils.DecodeCursor(cursor)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("(created_at, id) < ($%d, $%d)", argCount, argCount+1))
			args = append(args, cursorTime, cursorID)
			argCount += 2
		}
	}

	if len(conditions) > 0 {
		selectQuery += " AND " + strings.Join(conditions, " AND ")
	}

	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", argCount)
	args = append(args, limit+1)

	var products []models.Product
	err := r.db.SelectContext(ctx, &products, selectQuery, args...)
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(products) > limit {
		nextCursor = utils.EncodeCursor(products[limit].CreatedAt, products[limit].ID)
		products = products[:limit]
	}

	return products, nextCursor, nil
}

func (r *sqlCommandRepository) Create(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (id, name, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, p.ID, p.Name, p.Price, p.Stock).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *sqlCommandRepository) Update(ctx context.Context, p *models.Product) error {
	query := `
		UPDATE products
		SET name = $2, price = $3, stock = $4, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`
	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, p.ID, p.Name, p.Price, p.Stock).Scan(&p.UpdatedAt)
}

func (r *sqlCommandRepository) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	query := `
		UPDATE products
		SET stock = stock + $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	querier := database.GetQuerier(ctx, r.db)
	_, err := querier.ExecContext(ctx, query, id, quantityChange)
	return err
}

func (r *sqlCommandRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE products
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	querier := database.GetQuerier(ctx, r.db)
	_, err := querier.ExecContext(ctx, query, id)
	return err
}
