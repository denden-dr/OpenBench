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
	UpdateStatus(ctx context.Context, id string, newStatus string) error
	ClaimTicket(ctx context.Context, id string, technicianID string) error
	ListForBoard(ctx context.Context) ([]model.Ticket, error)
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
        SELECT id, device_type, brand, model, issue_description, status, diagnosis_fee, technician_id, created_at, updated_at 
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

func (r *sqlTicketRepository) UpdateStatus(ctx context.Context, id string, newStatus string) error {
	query := `UPDATE tickets SET status = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, newStatus, id)
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

func (r *sqlTicketRepository) ClaimTicket(ctx context.Context, id string, technicianID string) error {
	query := `
		UPDATE tickets
		SET status = $1, technician_id = $2
		WHERE id = $3 AND status = $4
	`
	result, err := r.db.ExecContext(ctx, query,
		model.StatusDiagnosing, technicianID, id, model.StatusReceived,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrClaimConflict
	}
	return nil
}

func (r *sqlTicketRepository) ListForBoard(ctx context.Context) ([]model.Ticket, error) {
	var tickets []model.Ticket
	query := `
		SELECT id, device_type, brand, model, status, created_at
		FROM tickets
		ORDER BY created_at DESC
	`
	if err := r.db.SelectContext(ctx, &tickets, query); err != nil {
		return nil, err
	}
	return tickets, nil
}
