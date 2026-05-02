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
        VALUES (:device_type, :brand, :model, :issue_description, :diagnosis_fee)
        RETURNING id, status, created_at, updated_at
    `
    // sqlx.NamedQueryContext doesn't support RETURNING into a struct easily with QueryRowx
    // We use NamedQuery then scan, or just use plain Exec if we don't need RETURNING 
    // but the model needs the ID and other defaults.
    
    rows, err := r.db.NamedQueryContext(ctx, query, ticket)
    if err != nil {
        return err
    }
    defer rows.Close()

    if rows.Next() {
        if err := rows.Scan(&ticket.ID, &ticket.Status, &ticket.CreatedAt, &ticket.UpdatedAt); err != nil {
            return err
        }
    }

    return rows.Err()
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
