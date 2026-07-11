package warranty

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type QueryRepository interface {
	FindWarrantyByID(ctx context.Context, id string) (*models.Warranty, error)
	FindWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error)
	FindClaimByID(ctx context.Context, id string) (*models.Claim, error)
	FindAllClaims(ctx context.Context, status string, search string, limit, offset int) ([]models.Claim, int, error)
}

type CommandRepository interface {
	CreateWarranty(ctx context.Context, w *models.Warranty) error
	UpdateWarrantyStatus(ctx context.Context, id string, status models.WarrantyStatus, notes *string) error
	CreateClaim(ctx context.Context, c *models.Claim) error
	UpdateClaim(ctx context.Context, c *models.Claim) error
	UpdateClaimEvaluation(ctx context.Context, claimID string, status models.ServiceTicketStatus, evalStatus models.ClaimEvaluationStatus, evalNotes *string) error
}

type sqlQueryRepository struct {
	db *sqlx.DB
}

type sqlCommandRepository struct {
	db *sqlx.DB
}

func NewQueryRepository(db *sqlx.DB) QueryRepository {
	return &sqlQueryRepository{db: db}
}

func NewCommandRepository(db *sqlx.DB) CommandRepository {
	return &sqlCommandRepository{db: db}
}

func (r *sqlCommandRepository) CreateWarranty(ctx context.Context, w *models.Warranty) error {
	query := `
		INSERT INTO warranties (id, ticket_id, start_date, end_date, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, w.ID, w.TicketID, w.StartDate, w.EndDate, w.Status, w.Notes).Scan(&w.CreatedAt, &w.UpdatedAt)
}

func (r *sqlQueryRepository) FindWarrantyByID(ctx context.Context, id string) (*models.Warranty, error) {
	query := `
		SELECT id, ticket_id, start_date, end_date, status, notes, created_at, updated_at
		FROM warranties
		WHERE id = $1
		LIMIT 1
	`
	var w models.Warranty
	err := r.db.GetContext(ctx, &w, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *sqlQueryRepository) FindWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error) {
	query := `
		SELECT id, ticket_id, start_date, end_date, status, notes, created_at, updated_at
		FROM warranties
		WHERE ticket_id = $1
		LIMIT 1
	`
	var w models.Warranty
	err := r.db.GetContext(ctx, &w, query, ticketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *sqlCommandRepository) UpdateWarrantyStatus(ctx context.Context, id string, status models.WarrantyStatus, notes *string) error {
	query := `
		UPDATE warranties
		SET status = $2, notes = $3, updated_at = NOW()
		WHERE id = $1
	`
	querier := database.GetQuerier(ctx, r.db)
	_, err := querier.ExecContext(ctx, query, id, status, notes)
	return err
}

func (r *sqlCommandRepository) CreateClaim(ctx context.Context, c *models.Claim) error {
	query := `
		INSERT INTO claims (id, claim_number, warranty_id, status, evaluation_status, issue_description, repair_action, notes, evaluation_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, c.ID, c.ClaimNumber, c.WarrantyID, c.Status, c.EvaluationStatus, c.IssueDescription, c.RepairAction, c.Notes, c.EvaluationNotes).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *sqlQueryRepository) FindClaimByID(ctx context.Context, id string) (*models.Claim, error) {
	query := `
		SELECT id, claim_number, warranty_id, status, evaluation_status, issue_description, repair_action, notes, evaluation_notes, created_at, updated_at
		FROM claims
		WHERE id = $1
		LIMIT 1
	`
	var c models.Claim
	err := r.db.GetContext(ctx, &c, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *sqlQueryRepository) FindAllClaims(ctx context.Context, status string, search string, limit, offset int) ([]models.Claim, int, error) {
	var selectQuery = `
		SELECT id, claim_number, warranty_id, status, evaluation_status, issue_description, repair_action, notes, evaluation_notes, created_at, updated_at
		FROM claims
	`
	var countQuery = `
		SELECT COUNT(*)
		FROM claims
	`

	var conditions []string
	var args []interface{}
	argCount := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, status)
		argCount++
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		conditions = append(conditions, fmt.Sprintf("(claim_number ILIKE $%d OR issue_description ILIKE $%d)", argCount, argCount))
		args = append(args, searchPattern)
		argCount++
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		selectQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination to select
	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	var claims []models.Claim
	err = r.db.SelectContext(ctx, &claims, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return claims, total, nil
}

func (r *sqlCommandRepository) UpdateClaim(ctx context.Context, c *models.Claim) error {
	query := `
		UPDATE claims
		SET 
			status = $2,
			issue_description = $3,
			repair_action = $4,
			notes = $5,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, c.ID, c.Status, c.IssueDescription, c.RepairAction, c.Notes).Scan(&c.UpdatedAt)
}

func (r *sqlCommandRepository) UpdateClaimEvaluation(ctx context.Context, claimID string, status models.ServiceTicketStatus, evalStatus models.ClaimEvaluationStatus, evalNotes *string) error {
	query := `
		UPDATE claims
		SET status = $2, evaluation_status = $3, evaluation_notes = $4, updated_at = NOW()
		WHERE id = $1
	`
	querier := database.GetQuerier(ctx, r.db)
	_, err := querier.ExecContext(ctx, query, claimID, status, evalStatus, evalNotes)
	return err
}
