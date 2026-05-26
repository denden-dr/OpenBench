package repository

import (
	"context"

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
	Update(ctx context.Context, ticket *model.Ticket) error
	UpdateTx(ctx context.Context, tx Transaction, ticket *model.Ticket) error
	List(ctx context.Context) ([]model.Ticket, error)
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
		SELECT id, customer_name, customer_gender, brand, model, issue,
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
		SET customer_name = $1, customer_gender = $2, brand = $3, model = $4, issue = $5,
		    additional_description = $6, accessories = $7, price = $8, status = $9,
		    payment_status = $10, warranty_days = $11, exit_date = $12,
		    is_warranty = $13, parent_ticket_id = $14
		WHERE id = $15
	`
	result, err := tx.ExecContext(ctx, query,
		ticket.CustomerName,
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
		INSERT INTO tickets (customer_name, customer_gender, brand, model, issue, additional_description, accessories, price, status, payment_status, warranty_days, is_warranty, parent_ticket_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, entry_date, status, payment_status
	`
	var err error
	if tx != nil {
		err = tx.QueryRowxContext(ctx, query,
			ticket.CustomerName,
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
		SELECT id, customer_name, customer_gender, brand, model, issue,
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

func (r *sqlTicketRepository) Update(ctx context.Context, ticket *model.Ticket) error {
	query := `
		UPDATE tickets
		SET customer_name = $1, customer_gender = $2, brand = $3, model = $4, issue = $5,
		    additional_description = $6, accessories = $7, price = $8, status = $9,
		    payment_status = $10, warranty_days = $11, exit_date = $12,
		    is_warranty = $13, parent_ticket_id = $14
		WHERE id = $15
	`
	result, err := r.db.ExecContext(ctx, query,
		ticket.CustomerName,
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
		SELECT id, customer_name, customer_gender, brand, model, issue,
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
