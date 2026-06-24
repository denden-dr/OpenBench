package sales

import (
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Sale struct {
	ID            string          `db:"id"`
	InvoiceNumber string          `db:"invoice_number"`
	Subtotal      decimal.Decimal `db:"subtotal"`
	Discount      decimal.Decimal `db:"discount"`
	Total         decimal.Decimal `db:"total"`
	PaymentMethod string          `db:"payment_method"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	Items         []SaleItem      `db:"-"`
}

type SaleItem struct {
	ID        string          `db:"id"`
	SaleID    string          `db:"sale_id"`
	ProductID string          `db:"product_id"`
	Name      string          `db:"name"` // Joined/derived from products
	Price     decimal.Decimal `db:"price"`
	Qty       int             `db:"qty"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func ToSaleAPI(s *Sale) api.Sale {
	itemsAPI := make([]api.SaleItem, len(s.Items))
	for i, item := range s.Items {
		price, _ := item.Price.Float64()
		itemsAPI[i] = api.SaleItem{
			ProductId: uuid.MustParse(item.ProductID),
			Name:      item.Name,
			Price:     float32(price),
			Qty:       item.Qty,
		}
	}

	subtotal, _ := s.Subtotal.Float64()
	discount, _ := s.Discount.Float64()
	total, _ := s.Total.Float64()

	return api.Sale{
		Id:            uuid.MustParse(s.ID),
		InvoiceNumber: s.InvoiceNumber,
		Subtotal:      float32(subtotal),
		Discount:      float32(discount),
		Total:         float32(total),
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
