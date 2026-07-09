package warranty

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateWarranty(ctx context.Context, w *models.Warranty) error
	FindWarrantyByID(ctx context.Context, id string) (*models.Warranty, error)
	FindWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error)
	UpdateWarrantyStatus(ctx context.Context, id string, status models.WarrantyStatus, notes *string) error

	CreateClaim(ctx context.Context, c *models.Claim) error
	FindClaimByID(ctx context.Context, id string) (*models.Claim, error)
	FindAllClaims(ctx context.Context, status string, search string, limit, offset int) ([]models.Claim, int, error)
	UpdateClaim(ctx context.Context, c *models.Claim) error
	EvaluateClaimTx(ctx context.Context, claimID string, evalStatus models.ClaimEvaluationStatus, evalNotes *string, repairStatus models.ServiceTicketStatus, isVoidWarranty bool, warrantyID string, warrantyNotes *string) error
}

type sqlRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) CreateWarranty(ctx context.Context, w *models.Warranty) error {
	query := `
		INSERT INTO warranties (id, ticket_id, start_date, end_date, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, w.ID, w.TicketID, w.StartDate, w.EndDate, w.Status, w.Notes).Scan(&w.CreatedAt, &w.UpdatedAt)
}

func (r *sqlRepository) FindWarrantyByID(ctx context.Context, id string) (*models.Warranty, error) {
	query := `
		SELECT id, ticket_id, start_date, end_date, status, notes, created_at, updated_at
		FROM warranties
		WHERE id = $1
		LIMIT 1
	`
	var w models.Warranty
	err := r.db.QueryRow(ctx, query, id).Scan(&w.ID, &w.TicketID, &w.StartDate, &w.EndDate, &w.Status, &w.Notes, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *sqlRepository) FindWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error) {
	query := `
		SELECT id, ticket_id, start_date, end_date, status, notes, created_at, updated_at
		FROM warranties
		WHERE ticket_id = $1
		LIMIT 1
	`
	var w models.Warranty
	err := r.db.QueryRow(ctx, query, ticketID).Scan(&w.ID, &w.TicketID, &w.StartDate, &w.EndDate, &w.Status, &w.Notes, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *sqlRepository) UpdateWarrantyStatus(ctx context.Context, id string, status models.WarrantyStatus, notes *string) error {
	query := `
		UPDATE warranties
		SET status = $2, notes = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id, status, notes)
	return err
}

func (r *sqlRepository) CreateClaim(ctx context.Context, c *models.Claim) error {
	query := `
		INSERT INTO claims (id, claim_number, warranty_id, status, evaluation_status, issue_description, repair_action, notes, evaluation_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, c.ID, c.ClaimNumber, c.WarrantyID, c.Status, c.EvaluationStatus, c.IssueDescription, c.RepairAction, c.Notes, c.EvaluationNotes).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *sqlRepository) FindClaimByID(ctx context.Context, id string) (*models.Claim, error) {
	query := `
		SELECT id, claim_number, warranty_id, status, evaluation_status, issue_description, repair_action, notes, evaluation_notes, created_at, updated_at
		FROM claims
		WHERE id = $1
		LIMIT 1
	`
	var c models.Claim
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.ClaimNumber, &c.WarrantyID, &c.Status, &c.EvaluationStatus, &c.IssueDescription, &c.RepairAction, &c.Notes, &c.EvaluationNotes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *sqlRepository) FindAllClaims(ctx context.Context, status string, search string, limit, offset int) ([]models.Claim, int, error) {
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
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination to select
	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var claims []models.Claim
	for rows.Next() {
		var c models.Claim
		err := rows.Scan(
			&c.ID, &c.ClaimNumber, &c.WarrantyID, &c.Status, &c.EvaluationStatus, &c.IssueDescription,
			&c.RepairAction, &c.Notes, &c.EvaluationNotes, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		claims = append(claims, c)
	}

	return claims, total, nil
}

func (r *sqlRepository) UpdateClaim(ctx context.Context, c *models.Claim) error {
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
	return r.db.QueryRow(ctx, query, c.ID, c.Status, c.IssueDescription, c.RepairAction, c.Notes).Scan(&c.UpdatedAt)
}

func (r *sqlRepository) EvaluateClaimTx(ctx context.Context, claimID string, evalStatus models.ClaimEvaluationStatus, evalNotes *string, repairStatus models.ServiceTicketStatus, isVoidWarranty bool, warrantyID string, warrantyNotes *string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Update claim status, evaluation status and notes
	queryClaim := `
		UPDATE claims
		SET status = $2, evaluation_status = $3, evaluation_notes = $4, updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, queryClaim, claimID, repairStatus, evalStatus, evalNotes)
	if err != nil {
		return err
	}

	// 2. If isVoidWarranty is true, update the associated warranty status to VOID and copy notes
	if isVoidWarranty {
		queryWarranty := `
			UPDATE warranties
			SET status = $2, notes = $3, updated_at = NOW()
			WHERE id = $1
		`
		_, err = tx.Exec(ctx, queryWarranty, warrantyID, models.WarrantyStatusVoid, warrantyNotes)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
