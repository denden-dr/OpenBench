package inventory

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

func (r *sqlQueryRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	query, args, err := r.psql.Select("id", "name", "price", "stock", "created_at", "updated_at").
		From("products").
		Where(squirrel.Eq{"id": id, "deleted_at": nil}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p models.Product
	err = r.db.GetContext(ctx, &p, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *sqlQueryRepository) FindAll(ctx context.Context, search string, limit int, cursor string) ([]models.Product, string, error) {
	queryBuilder := r.psql.Select("id", "name", "price", "stock", "created_at", "updated_at").
		From("products").
		Where(squirrel.Eq{"deleted_at": nil})

	if search != "" {
		searchPattern := "%" + search + "%"
		queryBuilder = queryBuilder.Where(squirrel.Expr("name ILIKE ?", searchPattern))
	}

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

	var products []models.Product
	err = r.db.SelectContext(ctx, &products, query, args...)
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
	query, args, err := r.psql.Insert("products").
		Columns("id", "name", "price", "stock", "created_at", "updated_at").
		Values(p.ID, p.Name, p.Price, p.Stock, squirrel.Expr("NOW()"), squirrel.Expr("NOW()")).
		Suffix("RETURNING created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, args...).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *sqlCommandRepository) Update(ctx context.Context, p *models.Product) error {
	query, args, err := r.psql.Update("products").
		Set("name", p.Name).
		Set("price", p.Price).
		Set("stock", p.Stock).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": p.ID, "deleted_at": nil}).
		Suffix("RETURNING updated_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, args...).Scan(&p.UpdatedAt)
}

func (r *sqlCommandRepository) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	query, args, err := r.psql.Update("products").
		Set("stock", squirrel.Expr("stock + ?", quantityChange)).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id, "deleted_at": nil}).
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	_, err = querier.ExecContext(ctx, query, args...)
	return err
}

func (r *sqlCommandRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.psql.Update("products").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id, "deleted_at": nil}).
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	_, err = querier.ExecContext(ctx, query, args...)
	return err
}
