package inventory

import (
	"errors"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInvalidInput      = errors.New("invalid input data")
	ErrProductReferenced = errors.New("cannot delete product, it is referenced in existing sales")
)

type Product struct {
	ID        string          `db:"id"`
	Name      string          `db:"name"`
	Category  string          `db:"category"`
	Stock     int             `db:"stock"`
	Price     decimal.Decimal `db:"price"`
	CostPrice decimal.Decimal `db:"cost_price"`
	MinStock  int             `db:"min_stock"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func ToProductAPI(p *Product) api.Product {
	price, _ := p.Price.Float64()
	costPrice, _ := p.CostPrice.Float64()
	return api.Product{
		Id:        uuid.MustParse(p.ID),
		Name:      p.Name,
		Category:  api.ProductCategory(p.Category),
		Stock:     p.Stock,
		Price:     float32(price),
		CostPrice: float32(costPrice),
		MinStock:  p.MinStock,
	}
}

func ToProductListAPI(products []*Product) []api.Product {
	res := make([]api.Product, len(products))
	for i, p := range products {
		res[i] = ToProductAPI(p)
	}
	return res
}
