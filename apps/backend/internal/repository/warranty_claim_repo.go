package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/denden-dr/openbench/apps/backend/internal/model"
	"github.com/jmoiron/sqlx"
)

// Compile-time check that *sqlx.Tx satisfies Transaction.
var _ Transaction = (*sqlx.Tx)(nil)

type WarrantyClaimRepository interface {
	BeginTx(ctx context.Context) (Transaction, error)
	Create(ctx context.Context, claim *model.WarrantyClaim) error
	GetByID(ctx context.Context, id string) (*model.WarrantyClaim, error)
	GetByIDForUpdateTx(ctx context.Context, tx Transaction, id string) (*model.WarrantyClaim, error)
	List(ctx context.Context, status string) ([]*model.WarrantyClaim, error)
	ListPaginated(ctx context.Context, status string, limit int, offset int) ([]*model.WarrantyClaim, error)
	CountPaginated(ctx context.Context, status string) (int64, error)
	UpdateTx(ctx context.Context, tx Transaction, claim *model.WarrantyClaim) error
	GetOpenClaimByTicketID(ctx context.Context, ticketID string) (*model.WarrantyClaim, error)
}

type sqlWarrantyClaimRepository struct {
	db *sqlx.DB
}

func NewWarrantyClaimRepository(db *sqlx.DB) WarrantyClaimRepository {
	return &sqlWarrantyClaimRepository{db: db}
}

func (r *sqlWarrantyClaimRepository) BeginTx(ctx context.Context) (Transaction, error) {
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

func (r *sqlWarrantyClaimRepository) GetByIDForUpdateTx(ctx context.Context, tx Transaction, id string) (*model.WarrantyClaim, error) {
	var claim model.WarrantyClaim
	query := `
		SELECT id, ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at, created_at, updated_at
		FROM warranty_claims
		WHERE id = $1
		FOR UPDATE
	`
	err := tx.GetContext(ctx, &claim, query, id)
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

func (r *sqlWarrantyClaimRepository) UpdateTx(ctx context.Context, tx Transaction, claim *model.WarrantyClaim) error {
	query := `
		UPDATE warranty_claims
		SET claim_ticket_id = $1, status = $2, void_reason = $3, inspected_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND status = 'waiting_inspection'
		RETURNING updated_at
	`
	err := tx.QueryRowxContext(ctx, query,
		claim.ClaimTicketID,
		claim.Status,
		claim.VoidReason,
		claim.InspectedAt,
		claim.ID,
	).Scan(&claim.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrConflict
		}
		return MapDatabaseError(err)
	}
	return nil
}

func (r *sqlWarrantyClaimRepository) ListPaginated(ctx context.Context, status string, limit int, offset int) ([]*model.WarrantyClaim, error) {
	var claims []*model.WarrantyClaim
	var err error
	if status != "" && status != "all" {
		query := `
			SELECT id, ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at, created_at, updated_at
			FROM warranty_claims
			WHERE status = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2 OFFSET $3
		`
		err = r.db.SelectContext(ctx, &claims, query, status, limit, offset)
	} else {
		query := `
			SELECT id, ticket_id, claim_ticket_id, issue, additional_description, status, void_reason, inspected_at, created_at, updated_at
			FROM warranty_claims
			ORDER BY created_at DESC, id DESC
			LIMIT $1 OFFSET $2
		`
		err = r.db.SelectContext(ctx, &claims, query, limit, offset)
	}
	if err != nil {
		return nil, MapDatabaseError(err)
	}
	return claims, nil
}

func (r *sqlWarrantyClaimRepository) CountPaginated(ctx context.Context, status string) (int64, error) {
	var count int64
	var err error
	if status != "" && status != "all" {
		query := `
			SELECT COUNT(*)
			FROM warranty_claims
			WHERE status = $1
		`
		err = r.db.GetContext(ctx, &count, query, status)
	} else {
		query := `
			SELECT COUNT(*)
			FROM warranty_claims
		`
		err = r.db.GetContext(ctx, &count, query)
	}
	if err != nil {
		return 0, MapDatabaseError(err)
	}
	return count, nil
}
