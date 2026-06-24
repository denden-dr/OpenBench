package sales

import (
	"context"
	"fmt"

	"github.com/denden-dr/openbench/apps/backend/internal/database"
	"github.com/jmoiron/sqlx"
)

type SalesRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, s *Sale) error
	CreateItem(ctx context.Context, tx *sqlx.Tx, item *SaleItem) error
	GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*Sale, error)
	List(ctx context.Context, tx *sqlx.Tx) ([]*Sale, error)
	GetNextInvoiceSequence(ctx context.Context, tx *sqlx.Tx, prefix string) (int, error)
	GetSaleItems(ctx context.Context, tx *sqlx.Tx, saleID string) ([]SaleItem, error)
}

type salesRepository struct {
	db *database.Database
}

func NewRepository(db *database.Database) SalesRepository {
	return &salesRepository{db: db}
}

func (r *salesRepository) Create(ctx context.Context, tx *sqlx.Tx, s *Sale) error {
	query := `
		INSERT INTO sales (
			id, invoice_number, subtotal, discount, total, payment_method
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, s.ID, s.InvoiceNumber, s.Subtotal, s.Discount, s.Total, s.PaymentMethod)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, s.ID, s.InvoiceNumber, s.Subtotal, s.Discount, s.Total, s.PaymentMethod)
	}
	return err
}

func (r *salesRepository) CreateItem(ctx context.Context, tx *sqlx.Tx, item *SaleItem) error {
	query := `
		INSERT INTO sale_items (
			id, sale_id, product_id, price, qty
		) VALUES ($1, $2, $3, $4, $5)
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, item.ID, item.SaleID, item.ProductID, item.Price, item.Qty)
	} else {
		_, err = r.db.DB.ExecContext(ctx, query, item.ID, item.SaleID, item.ProductID, item.Price, item.Qty)
	}
	return err
}

func (r *salesRepository) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*Sale, error) {
	var s Sale
	query := `
		SELECT id, invoice_number, subtotal, discount, total, payment_method, created_at, updated_at
		FROM sales WHERE id = $1
	`
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &s, query, id)
	} else {
		err = r.db.DB.GetContext(ctx, &s, query, id)
	}
	if err != nil {
		return nil, err
	}

	items, err := r.GetSaleItems(ctx, tx, s.ID)
	if err != nil {
		return nil, err
	}
	s.Items = items

	return &s, nil
}

func (r *salesRepository) List(ctx context.Context, tx *sqlx.Tx) ([]*Sale, error) {
	var sales []*Sale
	query := `
		SELECT id, invoice_number, subtotal, discount, total, payment_method, created_at, updated_at
		FROM sales ORDER BY created_at DESC
	`
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &sales, query)
	} else {
		err = r.db.DB.SelectContext(ctx, &sales, query)
	}
	if err != nil {
		return nil, err
	}

	if len(sales) == 0 {
		return sales, nil
	}

	var saleIDs []string
	saleMap := make(map[string]*Sale)
	for _, s := range sales {
		saleIDs = append(saleIDs, s.ID)
		s.Items = []SaleItem{}
		saleMap[s.ID] = s
	}

	queryItems, args, err := sqlx.In(`
		SELECT si.id, si.sale_id, si.product_id, p.name, si.price, si.qty, si.created_at, si.updated_at
		FROM sale_items si
		JOIN products p ON si.product_id = p.id
		WHERE si.sale_id IN (?)
	`, saleIDs)
	if err != nil {
		return nil, err
	}
	queryItems = r.db.DB.Rebind(queryItems)

	var items []SaleItem
	if tx != nil {
		err = tx.SelectContext(ctx, &items, queryItems, args...)
	} else {
		err = r.db.DB.SelectContext(ctx, &items, queryItems, args...)
	}
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if s, ok := saleMap[item.SaleID]; ok {
			s.Items = append(s.Items, item)
		}
	}

	return sales, nil
}

func (r *salesRepository) GetNextInvoiceSequence(ctx context.Context, tx *sqlx.Tx, prefix string) (int, error) {
	var maxNum int
	query := "SELECT COALESCE(MAX(CAST(SUBSTRING(invoice_number FROM '\\d+$') AS INTEGER)), 0) FROM sales WHERE invoice_number LIKE $1"
	likePattern := prefix + "%"
	var err error
	if tx != nil {
		// Acquire transaction lock to avoid race conditions
		lockQuery := "SELECT pg_advisory_xact_lock(hashtext($1))"
		_, err = tx.ExecContext(ctx, lockQuery, prefix)
		if err != nil {
			return 0, fmt.Errorf("failed to acquire advisory lock: %w", err)
		}
		err = tx.GetContext(ctx, &maxNum, query, likePattern)
	} else {
		err = r.db.DB.GetContext(ctx, &maxNum, query, likePattern)
	}
	return maxNum + 1, err
}

func (r *salesRepository) GetSaleItems(ctx context.Context, tx *sqlx.Tx, saleID string) ([]SaleItem, error) {
	var items []SaleItem
	query := `
		SELECT si.id, si.sale_id, si.product_id, p.name, si.price, si.qty, si.created_at, si.updated_at
		FROM sale_items si
		JOIN products p ON si.product_id = p.id
		WHERE si.sale_id = $1
	`
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &items, query, saleID)
	} else {
		err = r.db.DB.SelectContext(ctx, &items, query, saleID)
	}
	return items, err
}
