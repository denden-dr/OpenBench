package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/denden-dr/OpenBench/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BookingRepository defines the contract for booking data access.
type BookingRepository interface {
	Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error)
	FindByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*domain.Booking, error)
	UpdateStatus(ctx context.Context, id, userID uuid.UUID, currentStatus, newStatus domain.BookingStatus) error

	// Technician specific operations
	GetAvailableBookings(ctx context.Context) ([]*domain.Booking, error)
	GetBookingsByTechID(ctx context.Context, techID uuid.UUID) ([]*domain.Booking, error)
	UpdateTechnician(ctx context.Context, id, techID uuid.UUID) error
	UpdateDiagnosis(ctx context.Context, id, techID uuid.UUID, diagnosis string, cost float64, repairTime string) error
	UpdateTechStatus(ctx context.Context, id, techID uuid.UUID, currentStatus, newStatus domain.BookingStatus) error
}

type bookingRepository struct {
	db *sqlx.DB
}

// NewBookingRepository returns a new instance of the booking repository.
func NewBookingRepository(db *sqlx.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error) {
	query := `
		INSERT INTO bookings (user_id, device_name, issue_description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`
	var created domain.Booking
	err := r.db.QueryRowxContext(ctx, query, b.UserID, b.DeviceName, b.IssueDescription, b.Status).StructScan(&created)
	if err != nil {
		return nil, fmt.Errorf("creating booking: %w", err)
	}
	return &created, nil
}

func (r *bookingRepository) FindByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*domain.Booking, error) {
	query := `SELECT * FROM bookings WHERE id = $1 AND user_id = $2`
	var b domain.Booking
	err := r.db.GetContext(ctx, &b, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("finding booking: %w", err)
	}
	return &b, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id, userID uuid.UUID, currentStatus, newStatus domain.BookingStatus) error {
	query := `
		UPDATE bookings 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2 AND user_id = $3 AND status = $4
	`
	res, err := r.db.ExecContext(ctx, query, newStatus, id, userID, currentStatus)
	if err != nil {
		return fmt.Errorf("updating status: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrInvalidStateTransition
	}
	return nil
}

func (r *bookingRepository) GetAvailableBookings(ctx context.Context) ([]*domain.Booking, error) {
	query := `SELECT * FROM bookings WHERE technician_id IS NULL AND status = $1`
	var bookings []*domain.Booking
	err := r.db.SelectContext(ctx, &bookings, query, domain.StatusPendingDiagnosis)
	if err != nil {
		return nil, fmt.Errorf("getting available bookings: %w", err)
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingsByTechID(ctx context.Context, techID uuid.UUID) ([]*domain.Booking, error) {
	query := `SELECT * FROM bookings WHERE technician_id = $1`
	var bookings []*domain.Booking
	err := r.db.SelectContext(ctx, &bookings, query, techID)
	if err != nil {
		return nil, fmt.Errorf("getting tech bookings: %w", err)
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateTechnician(ctx context.Context, id, techID uuid.UUID) error {
	query := `
		UPDATE bookings 
		SET technician_id = $1, updated_at = NOW() 
		WHERE id = $2 AND technician_id IS NULL AND status = $3
	`
	res, err := r.db.ExecContext(ctx, query, techID, id, domain.StatusPendingDiagnosis)
	if err != nil {
		return fmt.Errorf("assigning technician: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrConflict
	}
	return nil
}

func (r *bookingRepository) UpdateDiagnosis(ctx context.Context, id, techID uuid.UUID, diagnosis string, cost float64, repairTime string) error {
	query := `
		UPDATE bookings 
		SET diagnosis_result = $1, estimated_cost = $2, estimated_repair_time = $3, status = $4, updated_at = NOW() 
		WHERE id = $5 AND technician_id = $6 AND status = $7
	`
	res, err := r.db.ExecContext(ctx, query, diagnosis, cost, repairTime, domain.StatusDiagnosisComplete, id, techID, domain.StatusPendingDiagnosis)
	if err != nil {
		return fmt.Errorf("updating diagnosis: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrInvalidStateTransition
	}
	return nil
}

func (r *bookingRepository) UpdateTechStatus(ctx context.Context, id, techID uuid.UUID, currentStatus, newStatus domain.BookingStatus) error {
	query := `
		UPDATE bookings 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2 AND technician_id = $3 AND status = $4
	`
	res, err := r.db.ExecContext(ctx, query, newStatus, id, techID, currentStatus)
	if err != nil {
		return fmt.Errorf("updating tech status: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrInvalidStateTransition
	}
	return nil
}
