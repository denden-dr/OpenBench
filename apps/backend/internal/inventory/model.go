package inventory

import (
	"errors"
	"time"

	"github.com/denden-dr/openbench/apps/backend/internal/pkg/api"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidInput    = errors.New("invalid input data")
)

type Product struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Category  string    `db:"category"`
	Stock     int       `db:"stock"`
	Price     float64   `db:"price"`
	CostPrice float64   `db:"cost_price"`
	MinStock  int       `db:"min_stock"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToProductAPI(p *Product) api.Product {
	return api.Product{
		Id:        uuid.MustParse(p.ID),
		Name:      p.Name,
		Category:  api.ProductCategory(p.Category),
		Stock:     p.Stock,
		Price:     float32(p.Price),
		CostPrice: float32(p.CostPrice),
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
