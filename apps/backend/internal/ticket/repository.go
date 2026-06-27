package ticket

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name=TicketRepository --output=mocks --outpkg=mocks --case=underscore
type TicketRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, t *Ticket) error
	GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*Ticket, error)
	GetByIDWithLock(ctx context.Context, tx *sqlx.Tx, id string) (*Ticket, error)
	GetByTicketNumber(ctx context.Context, tx *sqlx.Tx, ticketNumber string) (*Ticket, error)
	List(ctx context.Context, tx *sqlx.Tx) ([]*Ticket, error)
	Update(ctx context.Context, tx *sqlx.Tx, t *Ticket) error
	GetMaxTicketNumberByPrefix(ctx context.Context, tx *sqlx.Tx, prefix string) (string, error)

	CreateWarranty(ctx context.Context, tx *sqlx.Tx, w *Warranty) error
	GetWarrantyByTicketID(ctx context.Context, tx *sqlx.Tx, ticketID string) (*Warranty, error)
	ListWarranties(ctx context.Context, tx *sqlx.Tx) ([]*Warranty, error)
	DeleteWarrantyByTicketID(ctx context.Context, tx *sqlx.Tx, ticketID string) error
}

type ticketRepository struct {
	db *database.Database
}

// NewRepository creates a new postgres implementation of Repository
func NewRepository(db *database.Database) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, tx *sqlx.Tx, t *Ticket) error {
	query := `
		INSERT INTO tickets (
			id, ticket_number, customer_name, customer_phone, brand_phone, model_phone, 
			serial_number, damage_description, repair_action, cost, status, device_position, 
			payment_status, payment_method, warranty_duration_days, picked_up_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			t.ID, t.TicketNumber, t.CustomerName, t.CustomerPhone, t.BrandPhone, t.ModelPhone,
			t.SerialNumber, t.DamageDescription, t.RepairAction, t.Cost, t.Status, t.DevicePosition,
			t.PaymentStatus, t.PaymentMethod, t.WarrantyDurationDays, t.PickedUpAt)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query,
			t.ID, t.TicketNumber, t.CustomerName, t.CustomerPhone, t.BrandPhone, t.ModelPhone,
			t.SerialNumber, t.DamageDescription, t.RepairAction, t.Cost, t.Status, t.DevicePosition,
			t.PaymentStatus, t.PaymentMethod, t.WarrantyDurationDays, t.PickedUpAt)
	}
	if err != nil {
		return err
	}

	if t.Warranty != nil {
		upsertQuery := `
			INSERT INTO warranties (
				id, ticket_id, ticket_number, customer_name, device_info, start_date, end_date, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (ticket_id) DO UPDATE SET
				customer_name = EXCLUDED.customer_name,
				device_info = EXCLUDED.device_info,
				start_date = EXCLUDED.start_date,
				end_date = EXCLUDED.end_date,
				status = EXCLUDED.status,
				updated_at = CURRENT_TIMESTAMP
		`
		if tx != nil {
			_, err = tx.ExecContext(ctx, upsertQuery,
				t.Warranty.ID, t.Warranty.TicketID, t.Warranty.TicketNumber, t.Warranty.CustomerName,
				t.Warranty.DeviceInfo, t.Warranty.StartDate, t.Warranty.EndDate, t.Warranty.Status)
		} else {
			_, err = r.db.DB.ExecContext(ctx, upsertQuery,
				t.Warranty.ID, t.Warranty.TicketID, t.Warranty.TicketNumber, t.Warranty.CustomerName,
				t.Warranty.DeviceInfo, t.Warranty.StartDate, t.Warranty.EndDate, t.Warranty.Status)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ticketRepository) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*Ticket, error) {
	var t Ticket
	query := `
		SELECT id, ticket_number, customer_name, customer_phone, brand_phone, model_phone, 
		       serial_number, damage_description, repair_action, cost, status, device_position, 
		       payment_status, payment_method, warranty_duration_days, picked_up_at, 
		       created_at, updated_at 
		FROM tickets WHERE id = $1
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &t, query, id)
	} else {
		err = r.db.DB.GetContext(ctx, &t, query, id)
	}
	if err != nil {
		return nil, err
	}

	w, err := r.GetWarrantyByTicketID(ctx, tx, t.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	} else {
		t.Warranty = w
	}

	return &t, nil
}

func (r *ticketRepository) GetByIDWithLock(ctx context.Context, tx *sqlx.Tx, id string) (*Ticket, error) {
	var t Ticket
	query := `
		SELECT id, ticket_number, customer_name, customer_phone, brand_phone, model_phone, 
		       serial_number, damage_description, repair_action, cost, status, device_position, 
		       payment_status, payment_method, warranty_duration_days, picked_up_at, 
		       created_at, updated_at 
		FROM tickets WHERE id = $1 FOR UPDATE
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &t, query, id)
	} else {
		err = r.db.DB.GetContext(ctx, &t, query, id)
	}
	if err != nil {
		return nil, err
	}

	w, err := r.GetWarrantyByTicketID(ctx, tx, t.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	} else {
		t.Warranty = w
	}

	return &t, nil
}

func (r *ticketRepository) GetByTicketNumber(ctx context.Context, tx *sqlx.Tx, ticketNumber string) (*Ticket, error) {
	var t Ticket
	query := `
		SELECT id, ticket_number, customer_name, customer_phone, brand_phone, model_phone, 
		       serial_number, damage_description, repair_action, cost, status, device_position, 
		       payment_status, payment_method, warranty_duration_days, picked_up_at, 
		       created_at, updated_at 
		FROM tickets WHERE ticket_number = $1
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &t, query, ticketNumber)
	} else {
		err = r.db.DB.GetContext(ctx, &t, query, ticketNumber)
	}
	if err != nil {
		return nil, err
	}

	w, err := r.GetWarrantyByTicketID(ctx, tx, t.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	} else {
		t.Warranty = w
	}

	return &t, nil
}

func (r *ticketRepository) List(ctx context.Context, tx *sqlx.Tx) ([]*Ticket, error) {
	var tickets []*Ticket
	query := `
		SELECT id, ticket_number, customer_name, customer_phone, brand_phone, model_phone, 
		       serial_number, damage_description, repair_action, cost, status, device_position, 
		       payment_status, payment_method, warranty_duration_days, picked_up_at, 
		       created_at, updated_at 
		FROM tickets ORDER BY created_at DESC
	`
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &tickets, query)
	} else {
		err = r.db.DB.SelectContext(ctx, &tickets, query)
	}
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) Update(ctx context.Context, tx *sqlx.Tx, t *Ticket) error {
	query := `
		UPDATE tickets SET 
			customer_name = $1, customer_phone = $2, brand_phone = $3, model_phone = $4, 
			serial_number = $5, damage_description = $6, repair_action = $7, cost = $8, 
			status = $9, device_position = $10, payment_status = $11, payment_method = $12, 
			warranty_duration_days = $13, picked_up_at = $14, 
			updated_at = CURRENT_TIMESTAMP 
		WHERE id = $15
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			t.CustomerName, t.CustomerPhone, t.BrandPhone, t.ModelPhone,
			t.SerialNumber, t.DamageDescription, t.RepairAction, t.Cost,
			t.Status, t.DevicePosition, t.PaymentStatus, t.PaymentMethod,
			t.WarrantyDurationDays, t.PickedUpAt, t.ID)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query,
			t.CustomerName, t.CustomerPhone, t.BrandPhone, t.ModelPhone,
			t.SerialNumber, t.DamageDescription, t.RepairAction, t.Cost,
			t.Status, t.DevicePosition, t.PaymentStatus, t.PaymentMethod,
			t.WarrantyDurationDays, t.PickedUpAt, t.ID)
	}
	if err != nil {
		return err
	}

	if t.Warranty != nil {
		upsertQuery := `
			INSERT INTO warranties (
				id, ticket_id, ticket_number, customer_name, device_info, start_date, end_date, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (ticket_id) DO UPDATE SET
				customer_name = EXCLUDED.customer_name,
				device_info = EXCLUDED.device_info,
				start_date = EXCLUDED.start_date,
				end_date = EXCLUDED.end_date,
				status = EXCLUDED.status,
				updated_at = CURRENT_TIMESTAMP
		`
		if tx != nil {
			_, err = tx.ExecContext(ctx, upsertQuery,
				t.Warranty.ID, t.Warranty.TicketID, t.Warranty.TicketNumber, t.Warranty.CustomerName,
				t.Warranty.DeviceInfo, t.Warranty.StartDate, t.Warranty.EndDate, t.Warranty.Status)
		} else {
			_, err = r.db.DB.ExecContext(ctx, upsertQuery,
				t.Warranty.ID, t.Warranty.TicketID, t.Warranty.TicketNumber, t.Warranty.CustomerName,
				t.Warranty.DeviceInfo, t.Warranty.StartDate, t.Warranty.EndDate, t.Warranty.Status)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ticketRepository) GetMaxTicketNumberByPrefix(ctx context.Context, tx *sqlx.Tx, prefix string) (string, error) {
	var maxNum string
	query := "SELECT ticket_number FROM tickets WHERE ticket_number LIKE $1 ORDER BY created_at DESC, ticket_number DESC LIMIT 1"
	likePattern := prefix + "%"
	var err error
	if tx != nil {
		// Acquire transaction-level advisory lock on the prefix to prevent concurrent generation race
		lockQuery := "SELECT pg_advisory_xact_lock(hashtext($1))"
		_, err = tx.ExecContext(ctx, lockQuery, prefix)
		if err != nil {
			return "", fmt.Errorf("failed to acquire advisory lock for prefix: %w", err)
		}
		err = tx.GetContext(ctx, &maxNum, query, likePattern)
	} else {
		err = r.db.DB.GetContext(ctx, &maxNum, query, likePattern)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return maxNum, err
}

func (r *ticketRepository) CreateWarranty(ctx context.Context, tx *sqlx.Tx, w *Warranty) error {
	query := `
		INSERT INTO warranties (
			id, ticket_id, ticket_number, customer_name, device_info, start_date, end_date, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			w.ID, w.TicketID, w.TicketNumber, w.CustomerName, w.DeviceInfo, w.StartDate, w.EndDate, w.Status)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query,
			w.ID, w.TicketID, w.TicketNumber, w.CustomerName, w.DeviceInfo, w.StartDate, w.EndDate, w.Status)
	}
	return err
}

func (r *ticketRepository) GetWarrantyByTicketID(ctx context.Context, tx *sqlx.Tx, ticketID string) (*Warranty, error) {
	var w Warranty
	query := `
		SELECT id, ticket_id, ticket_number, customer_name, device_info, start_date, end_date, status, created_at, updated_at 
		FROM warranties WHERE ticket_id = $1
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &w, query, ticketID)
	} else {
		err = r.db.DB.GetContext(ctx, &w, query, ticketID)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *ticketRepository) ListWarranties(ctx context.Context, tx *sqlx.Tx) ([]*Warranty, error) {
	var warranties []*Warranty
	query := `
		SELECT id, ticket_id, ticket_number, customer_name, device_info, start_date, end_date, status, created_at, updated_at 
		FROM warranties ORDER BY created_at DESC
	`
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &warranties, query)
	} else {
		err = r.db.DB.SelectContext(ctx, &warranties, query)
	}
	if err != nil {
		return nil, err
	}
	return warranties, nil
}

func (r *ticketRepository) DeleteWarrantyByTicketID(ctx context.Context, tx *sqlx.Tx, ticketID string) error {
	query := `DELETE FROM warranties WHERE ticket_id = $1`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, ticketID)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, ticketID)
	}
	return err
}
