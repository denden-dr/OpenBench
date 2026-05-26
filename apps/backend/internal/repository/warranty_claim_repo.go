package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/jmoiron/sqlx"
)

type WarrantyClaimRepository interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	Create(ctx context.Context, claim *model.WarrantyClaim) error
	GetByID(ctx context.Context, id string) (*model.WarrantyClaim, error)
	List(ctx context.Context, status string) ([]*model.WarrantyClaim, error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, claim *model.WarrantyClaim) error
	GetOpenClaimByTicketID(ctx context.Context, ticketID string) (*model.WarrantyClaim, error)
}

type sqlWarrantyClaimRepository struct {
	db *sqlx.DB
}

func NewWarrantyClaimRepository(db *sqlx.DB) WarrantyClaimRepository {
	return &sqlWarrantyClaimRepository{db: db}
}

func (r *sqlWarrantyClaimRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *sqlWarrantyClaimRepository) Create(ctx context.Context, claim *model.WarrantyClaim) error {
	query := `
		INSERT INTO warranty_claims (ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowxContext(ctx, query,
		claim.TicketID,
		claim.ClaimTicketID,
		claim.Issue,
		claim.AdditionalDescription,
		claim.Status,
		claim.VoidReason,
		claim.InspectedAt,
	).Scan(&claim.ID, &claim.CreatedAt, &claim.UpdatedAt)
	return MapDatabaseError(err)
}

func (r *sqlWarrantyClaimRepository) GetByID(ctx context.Context, id string) (*model.WarrantyClaim, error) {
	var claim model.WarrantyClaim
	query := `
		SELECT id, ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at, created_at, updated_at
		FROM warranty_claims
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &claim, query, id)
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	return &claim, nil
}

func (r *sqlWarrantyClaimRepository) List(ctx context.Context, status string) ([]*model.WarrantyClaim, error) {
	var claims []*model.WarrantyClaim
	var err error
	if status != "" {
		query := `
			SELECT id, ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at, created_at, updated_at
			FROM warranty_claims
			WHERE status = $1
			ORDER BY created_at DESC
		`
		err = r.db.SelectContext(ctx, &claims, query, status)
	} else {
		query := `
			SELECT id, ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at, created_at, updated_at
			FROM warranty_claims
			ORDER BY created_at DESC
		`
		err = r.db.SelectContext(ctx, &claims, query)
	}
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	return claims, nil
}

func (r *sqlWarrantyClaimRepository) GetOpenClaimByTicketID(ctx context.Context, ticketID string) (*model.WarrantyClaim, error) {
	var claim model.WarrantyClaim
	query := `
		SELECT id, ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at, created_at, updated_at
		FROM warranty_claims
		WHERE ticket_id = $1 AND status = 'waiting_inspection'
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &claim, query, ticketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, MapDatabaseError(err)
	}
	return &claim, nil
}

func (r *sqlWarrantyClaimRepository) UpdateTx(ctx context.Context, tx *sqlx.Tx, claim *model.WarrantyClaim) error {
	query := `
		UPDATE warranty_claims
		SET claim_ticket_id = $1, status = $2, void_reason = $3, inspected_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND status = 'waiting_inspection'
	`
	result, err := tx.ExecContext(ctx, query,
		claim.ClaimTicketID,
		claim.Status,
		claim.VoidReason,
		claim.InspectedAt,
		claim.ID,
	)
	if err != nil {
		return MapDatabaseError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MapDatabaseError(err)
	}
	if rowsAffected == 0 {
		return ErrConflict
	}
	return nil
}
