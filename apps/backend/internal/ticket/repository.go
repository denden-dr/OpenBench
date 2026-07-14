package ticket

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
	FindAll(ctx context.Context, status string, search string, limit int, cursor string) ([]models.ServiceTicket, string, error)
	Search(ctx context.Context, req TicketSearchRequest) ([]models.ServiceTicket, string, error)
	FindByID(ctx context.Context, id string) (*models.ServiceTicket, error)
}

type CommandRepository interface {
	Create(ctx context.Context, ticket *models.ServiceTicket) error
	Update(ctx context.Context, ticket *models.ServiceTicket) error
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

func (r *sqlCommandRepository) Create(ctx context.Context, t *models.ServiceTicket) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := psql.Insert("service_tickets").
		Columns(
			"id", "ticket_number", "status", "customer_name", "customer_phone",
			"device_brand", "device_model", "device_passcode", "issue_description",
			"repair_action", "cost", "warranty_days", "notes", "created_at", "updated_at",
		).
		Values(
			t.ID, t.TicketNumber, t.Status, t.CustomerName, t.CustomerPhone,
			t.DeviceBrand, t.DeviceModel, t.DevicePasscode, t.IssueDescription,
			t.RepairAction, t.Cost, t.WarrantyDays, t.Notes, squirrel.Expr("NOW()"), squirrel.Expr("NOW()"),
		).
		Suffix("RETURNING created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	err = querier.QueryRowxContext(ctx, query, args...).Scan(&t.CreatedAt, &t.UpdatedAt)
	return err
}

func (r *sqlQueryRepository) FindAll(ctx context.Context, status string, search string, limit int, cursor string) ([]models.ServiceTicket, string, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := psql.Select(
		"id", "ticket_number", "status", "customer_name", "customer_phone",
		"device_brand", "device_model", "device_passcode", "issue_description",
		"repair_action", "cost", "warranty_days", "notes", "created_at", "updated_at",
	).From("service_tickets").Where(squirrel.Eq{"deleted_at": nil})

	if status != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": status})
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		queryBuilder = queryBuilder.Where(squirrel.Or{
			squirrel.Expr("ticket_number ILIKE ?", searchPattern),
			squirrel.Expr("customer_name ILIKE ?", searchPattern),
			squirrel.Expr("customer_phone ILIKE ?", searchPattern),
		})
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

	var tickets []models.ServiceTicket
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(tickets) > limit {
		nextCursor = utils.EncodeCursor(tickets[limit].CreatedAt, tickets[limit].ID)
		tickets = tickets[:limit]
	}

	return tickets, nextCursor, nil
}

func (r *sqlQueryRepository) FindByID(ctx context.Context, id string) (*models.ServiceTicket, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := psql.Select(
		"id", "ticket_number", "status", "customer_name", "customer_phone",
		"device_brand", "device_model", "device_passcode", "issue_description",
		"repair_action", "cost", "warranty_days", "notes", "created_at", "updated_at",
	).
		From("service_tickets").
		Where(squirrel.Eq{"id": id, "deleted_at": nil}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var t models.ServiceTicket
	err = r.db.GetContext(ctx, &t, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *sqlCommandRepository) Update(ctx context.Context, t *models.ServiceTicket) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := psql.Update("service_tickets").
		Set("ticket_number", t.TicketNumber).
		Set("status", t.Status).
		Set("customer_name", t.CustomerName).
		Set("customer_phone", t.CustomerPhone).
		Set("device_brand", t.DeviceBrand).
		Set("device_model", t.DeviceModel).
		Set("device_passcode", t.DevicePasscode).
		Set("issue_description", t.IssueDescription).
		Set("repair_action", t.RepairAction).
		Set("cost", t.Cost).
		Set("warranty_days", t.WarrantyDays).
		Set("notes", t.Notes).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": t.ID, "deleted_at": nil}).
		Suffix("RETURNING updated_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	err = querier.QueryRowxContext(ctx, query, args...).Scan(&t.UpdatedAt)
	return err
}

func (r *sqlQueryRepository) Search(ctx context.Context, req TicketSearchRequest) ([]models.ServiceTicket, string, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > utils.MaxLimit {
		limit = utils.MaxLimit
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := psql.Select(
		"id", "ticket_number", "status", "customer_name", "customer_phone",
		"device_brand", "device_model", "device_passcode", "issue_description",
		"repair_action", "cost", "warranty_days", "notes", "created_at", "updated_at",
	).From("service_tickets").Where(squirrel.Eq{"deleted_at": nil})

	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		queryBuilder = queryBuilder.Where(squirrel.Or{
			squirrel.Expr("id::text ILIKE ?", searchPattern),
			squirrel.Expr("ticket_number ILIKE ?", searchPattern),
			squirrel.Expr("customer_name ILIKE ?", searchPattern),
			squirrel.Expr("device_brand ILIKE ?", searchPattern),
			squirrel.Expr("device_model ILIKE ?", searchPattern),
		})
	}

	if req.ExactDate != "" {
		queryBuilder = queryBuilder.Where("created_at::date = ?", req.ExactDate)
	} else {
		if req.StartDate != "" {
			queryBuilder = queryBuilder.Where("created_at::date >= ?", req.StartDate)
		}
		if req.EndDate != "" {
			queryBuilder = queryBuilder.Where("created_at::date <= ?", req.EndDate)
		}
	}

	if req.IsActive != nil {
		if *req.IsActive {
			queryBuilder = queryBuilder.Where("status NOT IN ('COMPLETED', 'RETURNED')")
		} else {
			queryBuilder = queryBuilder.Where("status IN ('COMPLETED', 'RETURNED')")
		}
	}

	if req.Cursor != "" {
		cursorTime, cursorID, err := utils.DecodeCursor(req.Cursor)
		if err == nil {
			queryBuilder = queryBuilder.Where("(created_at, id) < (?, ?)", cursorTime, cursorID)
		}
	}

	queryBuilder = queryBuilder.OrderBy("created_at DESC", "id DESC").Limit(uint64(limit + 1))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, "", err
	}

	var tickets []models.ServiceTicket
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(tickets) > limit {
		nextCursor = utils.EncodeCursor(tickets[limit].CreatedAt, tickets[limit].ID)
		tickets = tickets[:limit]
	}

	return tickets, nextCursor, nil
}
