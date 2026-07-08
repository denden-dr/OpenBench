package ticket

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, ticket *models.ServiceTicket) error
	FindAll(ctx context.Context, status string, search string, limit, offset int) ([]models.ServiceTicket, int, error)
	FindByID(ctx context.Context, id string) (*models.ServiceTicket, error)
	Update(ctx context.Context, ticket *models.ServiceTicket) error
}

type sqlRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) Create(ctx context.Context, t *models.ServiceTicket) error {
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
	err := r.db.QueryRow(ctx, query,
		t.ID, t.TicketNumber, t.Status, t.CustomerName, t.CustomerPhone,
		t.DeviceBrand, t.DeviceModel, t.DevicePasscode, t.IssueDescription,
		t.RepairAction, t.Cost, t.WarrantyDays, t.Notes,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
	return err
}

func (r *sqlRepository) FindAll(ctx context.Context, status string, search string, limit, offset int) ([]models.ServiceTicket, int, error) {
	var selectQuery = `
		SELECT 
			id, ticket_number, status, customer_name, customer_phone, 
			device_brand, device_model, device_passcode, issue_description, 
			repair_action, cost, warranty_days, notes, created_at, updated_at
		FROM service_tickets
		WHERE deleted_at IS NULL
	`

	var countQuery = `
		SELECT COUNT(*)
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

	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		selectQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination to select
	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tickets []models.ServiceTicket
	for rows.Next() {
		var t models.ServiceTicket
		err := rows.Scan(
			&t.ID, &t.TicketNumber, &t.Status, &t.CustomerName, &t.CustomerPhone,
			&t.DeviceBrand, &t.DeviceModel, &t.DevicePasscode, &t.IssueDescription,
			&t.RepairAction, &t.Cost, &t.WarrantyDays, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		tickets = append(tickets, t)
	}

	return tickets, total, nil
}

func (r *sqlRepository) FindByID(ctx context.Context, id string) (*models.ServiceTicket, error) {
	query := `
		SELECT 
			id, ticket_number, status, customer_name, customer_phone, 
			device_brand, device_model, device_passcode, issue_description, 
			repair_action, cost, warranty_days, notes, created_at, updated_at
		FROM service_tickets
		WHERE id = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, query, id)

	var t models.ServiceTicket
	err := row.Scan(
		&t.ID, &t.TicketNumber, &t.Status, &t.CustomerName, &t.CustomerPhone,
		&t.DeviceBrand, &t.DeviceModel, &t.DevicePasscode, &t.IssueDescription,
		&t.RepairAction, &t.Cost, &t.WarrantyDays, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *sqlRepository) Update(ctx context.Context, t *models.ServiceTicket) error {
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
	err := r.db.QueryRow(ctx, query,
		t.ID, t.TicketNumber, t.Status, t.CustomerName, t.CustomerPhone,
		t.DeviceBrand, t.DeviceModel, t.DevicePasscode, t.IssueDescription,
		t.RepairAction, t.Cost, t.WarrantyDays, t.Notes,
	).Scan(&t.UpdatedAt)
	return err
}
