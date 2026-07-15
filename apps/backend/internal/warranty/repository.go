package warranty

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/denden-dr/OpenBench/apps/backend/internal/database"
	"github.com/denden-dr/OpenBench/apps/backend/internal/models"
	"github.com/denden-dr/OpenBench/apps/backend/internal/utils"
	"github.com/jmoiron/sqlx"
)

type QueryRepository interface {
	FindWarrantyByID(ctx context.Context, id string) (*models.Warranty, error)
	FindWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error)
	FindClaimByID(ctx context.Context, id string) (*models.Claim, error)
	FindAllClaims(ctx context.Context, status string, search string, limit int, cursor string) ([]models.Claim, string, error)
}

type CommandRepository interface {
	CreateWarranty(ctx context.Context, w *models.Warranty) error
	UpdateWarrantyStatus(ctx context.Context, id string, status models.WarrantyStatus, notes *string) error
	CreateClaim(ctx context.Context, c *models.Claim) error
	UpdateClaim(ctx context.Context, c *models.Claim) error
	UpdateClaimEvaluation(ctx context.Context, claimID string, status models.ServiceTicketStatus, evalStatus models.ClaimEvaluationStatus, evalNotes *string) error
}

type sqlQueryRepository struct {
	db   *sqlx.DB
	psql squirrel.StatementBuilderType
}

type sqlCommandRepository struct {
	db   *sqlx.DB
	psql squirrel.StatementBuilderType
}

func NewQueryRepository(db *sqlx.DB) QueryRepository {
	return &sqlQueryRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func NewCommandRepository(db *sqlx.DB) CommandRepository {
	return &sqlCommandRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *sqlCommandRepository) CreateWarranty(ctx context.Context, w *models.Warranty) error {
	query, args, err := r.psql.Insert("warranties").
		Columns("id", "ticket_id", "start_date", "end_date", "status", "notes", "created_at", "updated_at").
		Values(w.ID, w.TicketID, w.StartDate, w.EndDate, w.Status, w.Notes, squirrel.Expr("NOW()"), squirrel.Expr("NOW()")).
		Suffix("RETURNING created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, args...).Scan(&w.CreatedAt, &w.UpdatedAt)
}

func (r *sqlQueryRepository) FindWarrantyByID(ctx context.Context, id string) (*models.Warranty, error) {
	query, args, err := r.psql.Select("id", "ticket_id", "start_date", "end_date", "status", "notes", "created_at", "updated_at").
		From("warranties").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var w models.Warranty
	err = r.db.GetContext(ctx, &w, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *sqlQueryRepository) FindWarrantyByTicketID(ctx context.Context, ticketID string) (*models.Warranty, error) {
	query, args, err := r.psql.Select("id", "ticket_id", "start_date", "end_date", "status", "notes", "created_at", "updated_at").
		From("warranties").
		Where(squirrel.Eq{"ticket_id": ticketID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var w models.Warranty
	err = r.db.GetContext(ctx, &w, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *sqlCommandRepository) UpdateWarrantyStatus(ctx context.Context, id string, status models.WarrantyStatus, notes *string) error {
	query, args, err := r.psql.Update("warranties").
		Set("status", status).
		Set("notes", notes).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	_, err = querier.ExecContext(ctx, query, args...)
	return err
}

func (r *sqlCommandRepository) CreateClaim(ctx context.Context, c *models.Claim) error {
	query, args, err := r.psql.Insert("claims").
		Columns("id", "claim_number", "warranty_id", "status", "evaluation_status", "issue_description", "repair_action", "notes", "evaluation_notes", "created_at", "updated_at").
		Values(c.ID, c.ClaimNumber, c.WarrantyID, c.Status, c.EvaluationStatus, c.IssueDescription, c.RepairAction, c.Notes, c.EvaluationNotes, squirrel.Expr("NOW()"), squirrel.Expr("NOW()")).
		Suffix("RETURNING created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, args...).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *sqlQueryRepository) FindClaimByID(ctx context.Context, id string) (*models.Claim, error) {
	query, args, err := r.psql.Select("id", "claim_number", "warranty_id", "status", "evaluation_status", "issue_description", "repair_action", "notes", "evaluation_notes", "created_at", "updated_at").
		From("claims").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var c models.Claim
	err = r.db.GetContext(ctx, &c, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *sqlQueryRepository) FindAllClaims(ctx context.Context, status string, search string, limit int, cursor string) ([]models.Claim, string, error) {
	queryBuilder := r.psql.Select("id", "claim_number", "warranty_id", "status", "evaluation_status", "issue_description", "repair_action", "notes", "evaluation_notes", "created_at", "updated_at").
		From("claims")

	if status != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": status})
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		queryBuilder = queryBuilder.Where(squirrel.Or{
			squirrel.Expr("claim_number ILIKE ?", searchPattern),
			squirrel.Expr("issue_description ILIKE ?", searchPattern),
		})
	}

	if cursor != "" {
		cursorTime, cursorID, err := utils.DecodeCursor(cursor)
		if err == nil {
			queryBuilder = queryBuilder.Where("(created_at, id) < (?, ?)", cursorTime, cursorID)
		}
	}

	queryBuilder = queryBuilder.OrderBy("created_at DESC", "id DESC").Limit(uint64(limit + 1))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, "", err
	}

	var claims []models.Claim
	err = r.db.SelectContext(ctx, &claims, query, args...)
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(claims) > limit {
		nextCursor = utils.EncodeCursor(claims[limit].CreatedAt, claims[limit].ID)
		claims = claims[:limit]
	}

	return claims, nextCursor, nil
}

func (r *sqlCommandRepository) UpdateClaim(ctx context.Context, c *models.Claim) error {
	query, args, err := r.psql.Update("claims").
		Set("status", c.Status).
		Set("issue_description", c.IssueDescription).
		Set("repair_action", c.RepairAction).
		Set("notes", c.Notes).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": c.ID}).
		Suffix("RETURNING updated_at").
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	return querier.QueryRowxContext(ctx, query, args...).Scan(&c.UpdatedAt)
}

func (r *sqlCommandRepository) UpdateClaimEvaluation(ctx context.Context, claimID string, status models.ServiceTicketStatus, evalStatus models.ClaimEvaluationStatus, evalNotes *string) error {
	query, args, err := r.psql.Update("claims").
		Set("status", status).
		Set("evaluation_status", evalStatus).
		Set("evaluation_notes", evalNotes).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": claimID}).
		ToSql()
	if err != nil {
		return err
	}

	querier := database.GetQuerier(ctx, r.db)
	_, err = querier.ExecContext(ctx, query, args...)
	return err
}
