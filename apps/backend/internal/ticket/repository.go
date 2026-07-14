package ticket

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
	query := `
		INSERT INTO service_tickets (
			id, ticket_number, status, customer_name, customer_phone, 
			device_brand, device_model, device_passcode, issue_description, 
			repair_action, cost, warranty_days, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, 
			$10, $11, $12, $13, NOW(), NOW()
		)
		RETURNING created_at, updated_at
	`
	querier := database.GetQuerier(ctx, r.db)
	err := querier.QueryRowxContext(ctx, query,
		t.ID, t.TicketNumber, t.Status, t.CustomerName, t.CustomerPhone,
		t.DeviceBrand, t.DeviceModel, t.DevicePasscode, t.IssueDescription,
		t.RepairAction, t.Cost, t.WarrantyDays, t.Notes,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
	return err
}

func (r *sqlQueryRepository) FindAll(ctx context.Context, status string, search string, limit int, cursor string) ([]models.ServiceTicket, string, error) {
	var selectQuery = `
		SELECT 
			id, ticket_number, status, customer_name, customer_phone, 
			device_brand, device_model, device_passcode, issue_description, 
			repair_action, cost, warranty_days, notes, created_at, updated_at
		FROM service_tickets
		WHERE deleted_at IS NULL
	`

	var conditions []string
	var args []interface{}
	argCount := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, status)
		argCount++
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		conditions = append(conditions, fmt.Sprintf("(ticket_number ILIKE $%d OR customer_name ILIKE $%d OR customer_phone ILIKE $%d)", argCount, argCount, argCount))
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

	var tickets []models.ServiceTicket
	err := r.db.SelectContext(ctx, &tickets, selectQuery, args...)
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
	query := `
		SELECT 
			id, ticket_number, status, customer_name, customer_phone, 
			device_brand, device_model, device_passcode, issue_description, 
			repair_action, cost, warranty_days, notes, created_at, updated_at
		FROM service_tickets
		WHERE id = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var t models.ServiceTicket
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *sqlCommandRepository) Update(ctx context.Context, t *models.ServiceTicket) error {
	query := `
		UPDATE service_tickets
		SET 
			ticket_number = $2,
			status = $3,
			customer_name = $4,
			customer_phone = $5,
			device_brand = $6,
			device_model = $7,
			device_passcode = $8,
			issue_description = $9,
			repair_action = $10,
			cost = $11,
			warranty_days = $12,
			notes = $13,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`
	querier := database.GetQuerier(ctx, r.db)
	err := querier.QueryRowxContext(ctx, query,
		t.ID, t.TicketNumber, t.Status, t.CustomerName, t.CustomerPhone,
		t.DeviceBrand, t.DeviceModel, t.DevicePasscode, t.IssueDescription,
		t.RepairAction, t.Cost, t.WarrantyDays, t.Notes,
	).Scan(&t.UpdatedAt)
	return err
}

func (r *sqlQueryRepository) Search(ctx context.Context, req TicketSearchRequest) ([]models.ServiceTicket, string, error) {
	var selectQuery = `
		SELECT 
			id, ticket_number, status, customer_name, customer_phone, 
			device_brand, device_model, device_passcode, issue_description, 
			repair_action, cost, warranty_days, notes, created_at, updated_at
		FROM service_tickets
		WHERE deleted_at IS NULL
	`

	var conditions []string
	var args []interface{}
	argCount := 1

	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		conditions = append(conditions, fmt.Sprintf("(id::text ILIKE $%d OR ticket_number ILIKE $%d OR customer_name ILIKE $%d OR device_brand ILIKE $%d OR device_model ILIKE $%d)", argCount, argCount, argCount, argCount, argCount))
		args = append(args, searchPattern)
		argCount++
	}

	if req.ExactDate != "" {
		conditions = append(conditions, fmt.Sprintf("created_at::date = $%d", argCount))
		args = append(args, req.ExactDate)
		argCount++
	} else {
		if req.StartDate != "" {
			conditions = append(conditions, fmt.Sprintf("created_at::date >= $%d", argCount))
			args = append(args, req.StartDate)
			argCount++
		}
		if req.EndDate != "" {
			conditions = append(conditions, fmt.Sprintf("created_at::date <= $%d", argCount))
			args = append(args, req.EndDate)
			argCount++
		}
	}

	if req.IsActive != nil {
		if *req.IsActive {
			conditions = append(conditions, "status NOT IN ('COMPLETED', 'RETURNED')")
		} else {
			conditions = append(conditions, "status IN ('COMPLETED', 'RETURNED')")
		}
	}

	if req.Cursor != "" {
		cursorTime, cursorID, err := utils.DecodeCursor(req.Cursor)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("(created_at, id) < ($%d, $%d)", argCount, argCount+1))
			args = append(args, cursorTime, cursorID)
			argCount += 2
		}
	}

	if len(conditions) > 0 {
		selectQuery += " AND " + strings.Join(conditions, " AND ")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > utils.MaxLimit {
		limit = utils.MaxLimit
	}

	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", argCount)
	args = append(args, limit+1)

	var tickets []models.ServiceTicket
	err := r.db.SelectContext(ctx, &tickets, selectQuery, args...)
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
