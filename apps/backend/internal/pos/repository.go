package pos

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
	FindByID(ctx context.Context, id string) (*models.PosTransaction, error)
	FindAll(ctx context.Context, limit int, cursor string) ([]models.PosTransaction, string, error)
}

type CommandRepository interface {
	Create(ctx context.Context, t *models.PosTransaction) error
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

func (r *sqlQueryRepository) FindByID(ctx context.Context, id string) (*models.PosTransaction, error) {
	queryTx := `
		SELECT id, payment_method, total_amount, created_at
		FROM pos_transactions
		WHERE id = $1
		LIMIT 1
	`
	var t models.PosTransaction
	err := r.db.GetContext(ctx, &t, queryTx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	queryItems := `
		SELECT pti.id, pti.transaction_id, pti.product_id, pti.quantity, pti.price, COALESCE(p.name, 'Deleted Product') as product_name
		FROM pos_transaction_items pti
		LEFT JOIN products p ON pti.product_id = p.id
		WHERE pti.transaction_id = $1
	`
	err = r.db.SelectContext(ctx, &t.Items, queryItems, id)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *sqlQueryRepository) FindAll(ctx context.Context, limit int, cursor string) ([]models.PosTransaction, string, error) {
	var selectQuery = `
		SELECT id, payment_method, total_amount, created_at
		FROM pos_transactions
	`

	var conditions []string
	var args []interface{}
	argCount := 1

	if cursor != "" {
		cursorTime, cursorID, err := utils.DecodeCursor(cursor)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("(created_at, id) < ($%d, $%d)", argCount, argCount+1))
			args = append(args, cursorTime, cursorID)
			argCount += 2
		}
	}

	if len(conditions) > 0 {
		selectQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", argCount)
	args = append(args, limit+1)

	var transactions []models.PosTransaction
	err := r.db.SelectContext(ctx, &transactions, selectQuery, args...)
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(transactions) > limit {
		nextCursor = utils.EncodeCursor(transactions[limit].CreatedAt, transactions[limit].ID)
		transactions = transactions[:limit]
	}

	return transactions, nextCursor, nil
}

func (r *sqlCommandRepository) Create(ctx context.Context, t *models.PosTransaction) error {
	queryTx := `
		INSERT INTO pos_transactions (id, payment_method, total_amount, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING created_at
	`
	querier := database.GetQuerier(ctx, r.db)
	err := querier.QueryRowxContext(ctx, queryTx, t.ID, t.PaymentMethod, t.TotalAmount).Scan(&t.CreatedAt)
	if err != nil {
		return err
	}

	queryItem := `
		INSERT INTO pos_transaction_items (id, transaction_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
	`
	for i := range t.Items {
		item := &t.Items[i]
		_, err = querier.ExecContext(ctx, queryItem, item.ID, t.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}

	return nil
}
