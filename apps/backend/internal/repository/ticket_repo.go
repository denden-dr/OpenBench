package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/jmoiron/sqlx"
)

type TicketRepository interface {
	BeginTx(ctx context.Context) (Transaction, error)
	// Create persists a ticket. The caller must call ticket.PrepareForCreate() first.
	Create(ctx context.Context, ticket *model.Ticket) error
	// CreateTx is like Create but uses the given transaction.
	// The caller must call ticket.PrepareForCreate() first.
	CreateTx(ctx context.Context, tx Transaction, ticket *model.Ticket) error
	GetByID(ctx context.Context, id string) (*model.Ticket, error)
	GetByIDForUpdateTx(ctx context.Context, tx Transaction, id string) (*model.Ticket, error)
	GetByShortID(ctx context.Context, shortID string) ([]model.Ticket, error)
	GetByIDs(ctx context.Context, ids []string) ([]model.Ticket, error)
	Update(ctx context.Context, ticket *model.Ticket) error
	UpdateTx(ctx context.Context, tx Transaction, ticket *model.Ticket) error
	List(ctx context.Context) ([]model.Ticket, error)
	ListPaginated(ctx context.Context, search string, status string, limit int, offset int) ([]model.Ticket, error)
	CountPaginated(ctx context.Context, search string, status string) (int64, error)
	GetStatusCounts(ctx context.Context, search string) (map[string]int64, error)
	Delete(ctx context.Context, id string) error
}

type sqlTicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &sqlTicketRepository{db: db}
}

func (r *sqlTicketRepository) BeginTx(ctx context.Context) (Transaction, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *sqlTicketRepository) GetByIDForUpdateTx(ctx context.Context, tx Transaction, id string) (*model.Ticket, error) {
	var ticket model.Ticket
	query := `
		SELECT id, customer_name, customer_phone, customer_gender, brand, model, issue,
		       additional_description, accessories, price, status, payment_status,
		       warranty_days, entry_date, exit_date, is_warranty, parent_ticket_id
		FROM tickets
		WHERE id = $1
		FOR UPDATE
	`
	err := tx.GetContext(ctx, &ticket, query, id)
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	return &ticket, nil
}

func (r *sqlTicketRepository) UpdateTx(ctx context.Context, tx Transaction, ticket *model.Ticket) error {
	query := `
		UPDATE tickets
		SET customer_name = $1, customer_phone = $2, customer_gender = $3, brand = $4, model = $5, issue = $6,
		    additional_description = $7, accessories = $8, price = $9, status = $10,
		    payment_status = $11, warranty_days = $12, exit_date = $13,
		    is_warranty = $14, parent_ticket_id = $15
		WHERE id = $16
	`
	result, err := tx.ExecContext(ctx, query,
		ticket.CustomerName,
		ticket.CustomerPhone,
		ticket.CustomerGender,
		ticket.Brand,
		ticket.Model,
		ticket.Issue,
		ticket.AdditionalDescription,
		ticket.Accessories,
		ticket.Price,
		ticket.Status,
		ticket.PaymentStatus,
		ticket.WarrantyDays,
		ticket.ExitDate,
		ticket.IsWarranty,
		ticket.ParentTicketID,
		ticket.ID,
	)
	if err != nil {
		return MapDatabaseError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return MapDatabaseError(err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqlTicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	return r.CreateTx(ctx, nil, ticket)
}

func (r *sqlTicketRepository) CreateTx(ctx context.Context, tx Transaction, ticket *model.Ticket) error {
	query := `
		INSERT INTO tickets (customer_name, customer_phone, customer_gender, brand, model, issue, additional_description, accessories, price, status, payment_status, warranty_days, is_warranty, parent_ticket_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, entry_date, status, payment_status
	`
	var err error
	if tx != nil {
		err = tx.QueryRowxContext(ctx, query,
			ticket.CustomerName,
			ticket.CustomerPhone,
			ticket.CustomerGender,
			ticket.Brand,
			ticket.Model,
			ticket.Issue,
			ticket.AdditionalDescription,
			ticket.Accessories,
			ticket.Price,
			ticket.Status,
			ticket.PaymentStatus,
			ticket.WarrantyDays,
			ticket.IsWarranty,
			ticket.ParentTicketID,
		).Scan(&ticket.ID, &ticket.EntryDate, &ticket.Status, &ticket.PaymentStatus)
	} else {
		err = r.db.QueryRowxContext(ctx, query,
			ticket.CustomerName,
			ticket.CustomerPhone,
			ticket.CustomerGender,
			ticket.Brand,
			ticket.Model,
			ticket.Issue,
			ticket.AdditionalDescription,
			ticket.Accessories,
			ticket.Price,
			ticket.Status,
			ticket.PaymentStatus,
			ticket.WarrantyDays,
			ticket.IsWarranty,
			ticket.ParentTicketID,
		).Scan(&ticket.ID, &ticket.EntryDate, &ticket.Status, &ticket.PaymentStatus)
	}
	return MapDatabaseError(err)
}

func (r *sqlTicketRepository) GetByID(ctx context.Context, id string) (*model.Ticket, error) {
	var ticket model.Ticket
	query := `
		SELECT id, customer_name, customer_phone, customer_gender, brand, model, issue,
		       additional_description, accessories, price, status, payment_status,
		       warranty_days, entry_date, exit_date, is_warranty, parent_ticket_id
		FROM tickets
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	return &ticket, nil
}

func (r *sqlTicketRepository) GetByShortID(ctx context.Context, shortID string) ([]model.Ticket, error) {
	var tickets []model.Ticket
	query := `
		SELECT id, customer_name, customer_phone, customer_gender, brand, model, issue,
		       additional_description, accessories, price, status, payment_status,
		       warranty_days, entry_date, exit_date, is_warranty, parent_ticket_id
		FROM tickets
		WHERE left(id::text, 8) = lower($1)
		ORDER BY entry_date DESC
	`
	if err := r.db.SelectContext(ctx, &tickets, query, shortID); err != nil {
		return nil, MapDatabaseError(err)
	}
	return tickets, nil
}

func (r *sqlTicketRepository) Update(ctx context.Context, ticket *model.Ticket) error {
	query := `
		UPDATE tickets
		SET customer_name = $1, customer_phone = $2, customer_gender = $3, brand = $4, model = $5, issue = $6,
		    additional_description = $7, accessories = $8, price = $9, status = $10,
		    payment_status = $11, warranty_days = $12, exit_date = $13,
		    is_warranty = $14, parent_ticket_id = $15
		WHERE id = $16
	`
	result, err := r.db.ExecContext(ctx, query,
		ticket.CustomerName,
		ticket.CustomerPhone,
		ticket.CustomerGender,
		ticket.Brand,
		ticket.Model,
		ticket.Issue,
		ticket.AdditionalDescription,
		ticket.Accessories,
		ticket.Price,
		ticket.Status,
		ticket.PaymentStatus,
		ticket.WarrantyDays,
		ticket.ExitDate,
		ticket.IsWarranty,
		ticket.ParentTicketID,
		ticket.ID,
	)
	if err != nil {
		return MapDatabaseError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return MapDatabaseError(err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqlTicketRepository) List(ctx context.Context) ([]model.Ticket, error) {
	var tickets []model.Ticket
	query := `
		SELECT id, customer_name, customer_phone, customer_gender, brand, model, issue,
		       additional_description, accessories, price, status, payment_status,
		       warranty_days, entry_date, exit_date, is_warranty, parent_ticket_id
		FROM tickets
		ORDER BY entry_date DESC
	`
	if err := r.db.SelectContext(ctx, &tickets, query); err != nil {
		return nil, MapDatabaseError(err)
	}
	return tickets, nil
}

func (r *sqlTicketRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tickets WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return MapDatabaseError(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return MapDatabaseError(err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqlTicketRepository) GetByIDs(ctx context.Context, ids []string) ([]model.Ticket, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var tickets []model.Ticket
	query, args, err := sqlx.In(`
		SELECT id, customer_name, customer_phone, customer_gender, brand, model, issue,
		       additional_description, accessories, price, status, payment_status,
		       warranty_days, entry_date, exit_date, is_warranty, parent_ticket_id
		FROM tickets
		WHERE id IN (?)
	`, ids)
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	return tickets, nil
}

func (r *sqlTicketRepository) ListPaginated(ctx context.Context, search string, status string, limit int, offset int) ([]model.Ticket, error) {
	var tickets []model.Ticket
	query := `
		SELECT id, customer_name, customer_phone, customer_gender, brand, model, issue,
		       additional_description, accessories, price, status, payment_status,
		       warranty_days, entry_date, exit_date, is_warranty, parent_ticket_id
		FROM tickets
		WHERE 1=1
	`
	args := []interface{}{}
	paramIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND lower(
			COALESCE(id::text, '') || ' ' ||
			COALESCE(customer_name, '') || ' ' ||
			COALESCE(customer_phone, '') || ' ' ||
			COALESCE(brand, '') || ' ' ||
			COALESCE(model, '') || ' ' ||
			COALESCE(issue, '')
		) LIKE $%d`, paramIdx)
		args = append(args, "%"+strings.ToLower(search)+"%")
		paramIdx++
	}

	if status == "all" || status == "" {
		query += ` AND status != 'picked_up'`
	} else {
		query += fmt.Sprintf(` AND status = $%d`, paramIdx)
		args = append(args, status)
		paramIdx++
	}

	query += fmt.Sprintf(` ORDER BY entry_date DESC, id DESC LIMIT $%d OFFSET $%d`, paramIdx, paramIdx+1)
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	return tickets, nil
}

func (r *sqlTicketRepository) CountPaginated(ctx context.Context, search string, status string) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM tickets
		WHERE 1=1
	`
	args := []interface{}{}
	paramIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND lower(
			COALESCE(id::text, '') || ' ' ||
			COALESCE(customer_name, '') || ' ' ||
			COALESCE(customer_phone, '') || ' ' ||
			COALESCE(brand, '') || ' ' ||
			COALESCE(model, '') || ' ' ||
			COALESCE(issue, '')
		) LIKE $%d`, paramIdx)
		args = append(args, "%"+strings.ToLower(search)+"%")
		paramIdx++
	}

	if status == "all" || status == "" {
		query += ` AND status != 'picked_up'`
	} else {
		query += fmt.Sprintf(` AND status = $%d`, paramIdx)
		args = append(args, status)
		paramIdx++
	}

	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, MapDatabaseError(err)
	}
	return count, nil
}

func (r *sqlTicketRepository) GetStatusCounts(ctx context.Context, search string) (map[string]int64, error) {
	query := `
		SELECT status, COUNT(*)
		FROM tickets
		WHERE 1=1
	`
	args := []interface{}{}
	paramIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND lower(
			COALESCE(id::text, '') || ' ' ||
			COALESCE(customer_name, '') || ' ' ||
			COALESCE(customer_phone, '') || ' ' ||
			COALESCE(brand, '') || ' ' ||
			COALESCE(model, '') || ' ' ||
			COALESCE(issue, '')
		) LIKE $%d`, paramIdx)
		args = append(args, "%"+strings.ToLower(search)+"%")
		paramIdx++
	}

	query += ` GROUP BY status`

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	defer rows.Close()

	counts := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, MapDatabaseError(err)
		}
		counts[status] = count
	}
	if err := rows.Err(); err != nil {
		return nil, MapDatabaseError(err)
	}
	return counts, nil
}
