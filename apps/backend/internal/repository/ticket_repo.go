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
}

type sqlTicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &sqlTicketRepository{db: db}
}

func (r *sqlTicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	query := `
        INSERT INTO tickets (device_type, brand, model, issue_description, diagnosis_fee)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, status, created_at, updated_at
    `

	err := r.db.QueryRowxContext(ctx, query,
		ticket.DeviceType,
		ticket.Brand,
		ticket.Model,
		ticket.IssueDescription,
		ticket.DiagnosisFee,
	).Scan(&ticket.ID, &ticket.Status, &ticket.CreatedAt, &ticket.UpdatedAt)

	return err
}

func (r *sqlTicketRepository) GetByID(ctx context.Context, id string) (*model.Ticket, error) {
	var ticket model.Ticket
	query := `
        SELECT id, device_type, brand, model, issue_description, status, diagnosis_fee, created_at, updated_at 
        FROM tickets 
        WHERE id = $1
    `

	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &ticket, nil
}
