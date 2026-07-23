package dashboard

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetActiveTicketsCount(ctx context.Context) (int, error)
	GetPendingDiagnosesCount(ctx context.Context) (int, error)
	GetSalesToday(ctx context.Context) (float64, error)
	GetActiveWarrantiesCount(ctx context.Context) (int, error)
	GetRecentTickets(ctx context.Context) ([]RecentTicket, error)
}

type sqlRepository struct {
	db   *sqlx.DB
	psql squirrel.StatementBuilderType
}

func NewRepository(db *sqlx.DB) Repository {
	return &sqlRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *sqlRepository) GetActiveTicketsCount(ctx context.Context) (int, error) {
	query, args, err := r.psql.
		Select("COUNT(*)").
		From("service_tickets").
		Where("status NOT IN ('COMPLETED', 'RETURNED')").
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *sqlRepository) GetPendingDiagnosesCount(ctx context.Context) (int, error) {
	query, args, err := r.psql.
		Select("COUNT(*)").
		From("service_tickets").
		Where(squirrel.Eq{"status": "RECEIVED"}).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *sqlRepository) GetSalesToday(ctx context.Context) (float64, error) {
	query, args, err := r.psql.
		Select("COALESCE(SUM(total_amount), 0)").
		From("pos_transactions").
		Where("created_at::date = CURRENT_DATE").
		ToSql()
	if err != nil {
		return 0, err
	}

	var total float64
	err = r.db.GetContext(ctx, &total, query, args...)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *sqlRepository) GetActiveWarrantiesCount(ctx context.Context) (int, error) {
	query, args, err := r.psql.
		Select("COUNT(*)").
		From("warranties").
		Where(squirrel.Eq{"status": "ACTIVE"}).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *sqlRepository) GetRecentTickets(ctx context.Context) ([]RecentTicket, error) {
	query, args, err := r.psql.
		Select("id", "ticket_number", "status", "customer_name", "device_brand", "device_model", "cost", "created_at").
		From("service_tickets").
		OrderBy("updated_at DESC").
		Limit(4).
		ToSql()
	if err != nil {
		return nil, err
	}

	var tickets []RecentTicket
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}
