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

// TechnicianRepository defines the contract for technician data access.
type TechnicianRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Technician, error)
}

// technicianRepository is the private implementation of TechnicianRepository.
type technicianRepository struct {
	db *sqlx.DB
}

// NewTechnicianRepository returns a new instance of the technician repository.
func NewTechnicianRepository(db *sqlx.DB) TechnicianRepository {
	return &technicianRepository{db: db}
}

func (r *technicianRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Technician, error) {
	query := `SELECT * FROM technicians WHERE user_id = $1`
	var t domain.Technician
	err := r.db.GetContext(ctx, &t, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("finding technician by user id: %w", err)
	}
	return &t, nil
}
