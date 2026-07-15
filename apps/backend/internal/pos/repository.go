package pos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
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
	db   *sqlx.DB
	psql squirrel.StatementBuilderType
}

type sqlCommandRepository struct {
	db   *sqlx.DB
	psql squirrel.StatementBuilderType
}

func NewQueryRepository(db *sqlx.DB) QueryRepository {
	return &sqlQueryRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func NewCommandRepository(db *sqlx.DB) CommandRepository {
	return &sqlCommandRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *sqlQueryRepository) FindByID(ctx context.Context, id string) (*models.PosTransaction, error) {
	queryTx, argsTx, err := r.psql.Select("id", "payment_method", "total_amount", "created_at").
		From("pos_transactions").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var t models.PosTransaction
	err = r.db.GetContext(ctx, &t, queryTx, argsTx...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	queryItems, argsItems, err := r.psql.Select(
		"pti.id", "pti.transaction_id", "pti.product_id", "pti.quantity", "pti.price",
		"COALESCE(p.name, 'Deleted Product') as product_name",
	).
		From("pos_transaction_items pti").
		LeftJoin("products p ON pti.product_id = p.id").
		Where(squirrel.Eq{"pti.transaction_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.SelectContext(ctx, &t.Items, queryItems, argsItems...)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *sqlQueryRepository) FindAll(ctx context.Context, limit int, cursor string) ([]models.PosTransaction, string, error) {
	queryBuilder := r.psql.Select("id", "payment_method", "total_amount", "created_at").
		From("pos_transactions")

	if cursor != "" {
		cursorTime, cursorID, err := utils.DecodeCursor(cursor)
		if err == nil {
			queryBuilder = queryBuilder.Where("(created_at, id) < (?, ?)", cursorTime, cursorID)
		}
	}

	queryBuilder = queryBuilder.OrderBy("created_at DESC", "id DESC").Limit(uint64(limit + 1))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, "", err
	}

	var transactions []models.PosTransaction
	err = r.db.SelectContext(ctx, &transactions, query, args...)
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
	queryTx, argsTx, err := r.psql.Insert("pos_transactions").
		Columns("id", "payment_method", "total_amount", "created_at").
		Values(t.ID, t.PaymentMethod, t.TotalAmount, squirrel.Expr("NOW()")).
		Suffix("RETURNING created_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	err = querier.QueryRowxContext(ctx, queryTx, argsTx...).Scan(&t.CreatedAt)
	if err != nil {
		return err
	}

	if len(t.Items) > 0 {
		insertBuilder := r.psql.Insert("pos_transaction_items").
			Columns("id", "transaction_id", "product_id", "quantity", "price")

		for i := range t.Items {
			item := &t.Items[i]
			insertBuilder = insertBuilder.Values(item.ID, t.ID, item.ProductID, item.Quantity, item.Price)
		}

		queryItem, argsItem, err := insertBuilder.ToSql()
		if err != nil {
			return err
		}

		_, err = querier.ExecContext(ctx, queryItem, argsItem...)
		if err != nil {
			return err
		}
	}

	return nil
}
