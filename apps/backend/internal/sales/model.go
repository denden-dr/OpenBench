package sales

import (
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
)

type Sale struct {
	ID            string     `db:"id"`
	InvoiceNumber string     `db:"invoice_number"`
	Subtotal      float64    `db:"subtotal"`
	Discount      float64    `db:"discount"`
	Total         float64    `db:"total"`
	PaymentMethod string     `db:"payment_method"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	Items         []SaleItem `db:"-"`
}

type SaleItem struct {
	ID        string    `db:"id"`
	SaleID    string    `db:"sale_id"`
	ProductID string    `db:"product_id"`
	Name      string    `db:"name"` // Joined/derived from products
	Price     float64   `db:"price"`
	Qty       int       `db:"qty"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToSaleAPI(s *Sale) api.Sale {
	itemsAPI := make([]api.SaleItem, len(s.Items))
	for i, item := range s.Items {
		itemsAPI[i] = api.SaleItem{
			ProductId: uuid.MustParse(item.ProductID),
			Name:      item.Name,
			Price:     float32(item.Price),
			Qty:       item.Qty,
		}
	}

	return api.Sale{
		Id:            uuid.MustParse(s.ID),
		InvoiceNumber: s.InvoiceNumber,
		Subtotal:      float32(s.Subtotal),
		Discount:      float32(s.Discount),
		Total:         float32(s.Total),
		PaymentMethod: api.SalePaymentMethod(s.PaymentMethod),
		CreatedAt:     s.CreatedAt,
		Items:         itemsAPI,
	}
}

func ToSaleListAPI(sales []*Sale) []api.Sale {
	res := make([]api.Sale, len(sales))
	for i, s := range sales {
		res[i] = ToSaleAPI(s)
	}
	return res
}
