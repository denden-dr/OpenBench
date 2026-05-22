package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/jmoiron/sqlx"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket *model.Ticket) error
	GetByID(ctx context.Context, id string) (*model.Ticket, error)
	Update(ctx context.Context, ticket *model.Ticket) error
	List(ctx context.Context) ([]model.Ticket, error)
	Delete(ctx context.Context, id string) error
}

type sqlTicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &sqlTicketRepository{db: db}
}

func (r *sqlTicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	query := `
		INSERT INTO tickets (customer_name, customer_gender, brand, model, issue, additional_description, accessories, price, status, payment_status, warranty_days)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, entry_date, status, payment_status
	`
	return r.db.QueryRowxContext(ctx, query,
		ticket.CustomerName,
		ticket.CustomerGender,
		ticket.Brand,
		ticket.Model,
		ticket.Issue,
		ticket.AdditionalDescription,
		ticket.Accessories,
		ticket.Price,
		"service_in", // Default initial status
		"unpaid",     // Default initial payment status
		ticket.WarrantyDays,
	).Scan(&ticket.ID, &ticket.EntryDate, &ticket.Status, &ticket.PaymentStatus)
}

func (r *sqlTicketRepository) GetByID(ctx context.Context, id string) (*model.Ticket, error) {
	var ticket model.Ticket
	query := `SELECT * FROM tickets WHERE id = $1`
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *sqlTicketRepository) Update(ctx context.Context, ticket *model.Ticket) error {
	query := `
		UPDATE tickets
		SET customer_name = $1, customer_gender = $2, brand = $3, model = $4, issue = $5,
		    additional_description = $6, accessories = $7, price = $8, status = $9,
		    payment_status = $10, warranty_days = $11, exit_date = $12, warranty_expiry_date = $13
		WHERE id = $14
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
		ticket.WarrantyExpiryDate,
		ticket.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *sqlTicketRepository) List(ctx context.Context) ([]model.Ticket, error) {
	var tickets []model.Ticket
	query := `SELECT * FROM tickets ORDER BY entry_date DESC`
	if err := r.db.SelectContext(ctx, &tickets, query); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *sqlTicketRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tickets WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
