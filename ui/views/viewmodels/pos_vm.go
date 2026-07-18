package viewmodels

import (
	"time"

	"github.com/denden-dr/OpenBench/internal/models"
	"github.com/denden-dr/OpenBench/internal/utils"
)

// ProductVM is the view model for a product item in the POS/Inventory UI.
// Templates should use this instead of models.Product directly.
type ProductVM struct {
	// Raw values needed for Alpine.js client-side logic
	ID    string
	Name  string
	Price int64
	Stock int

	// Pre-formatted values for direct rendering in templates
	PriceFormatted  string // e.g. "150.000"
	StockText       string // e.g. "12" or "0"
	StockLabel      string // e.g. "Available", "Low Stock", "Out of Stock"
	StockBadgeClass string // Tailwind classes for the badge color
	IsOutOfStock    bool
	IDShort         string // first 8 chars of UUID
}

// TransactionVM is the view model for a POS transaction in the history UI.
// Templates should use this instead of models.PosTransaction directly.
type TransactionVM struct {
	IDShort              string // e.g. "a1b2c3d4..."
	PaymentMethodLabel   string // "CASH" or "QRIS"
	IsCash               bool
	TotalAmountFormatted string // e.g. "Rp 150.000"
	CreatedAtFormatted   string // e.g. "17 Jul 2026, 15:04 WIB"
}

// NewProductVM maps a models.Product to a ProductVM with pre-computed display values.
func NewProductVM(p models.Product) ProductVM {
	stockText := utils.IntToString(p.Stock)
	label, badgeClass := stockBadge(p.Stock)

	return ProductVM{
		ID:              p.ID,
		Name:            p.Name,
		Price:           p.Price,
		Stock:           p.Stock,
		PriceFormatted:  utils.FormatCurrency(p.Price),
		StockText:       stockText,
		StockLabel:      label,
		StockBadgeClass: badgeClass,
		IsOutOfStock:    p.Stock <= 0,
		IDShort:         p.ID[:8],
	}
}

// NewProductVMs maps a slice of models.Product to a slice of ProductVM.
func NewProductVMs(products []models.Product) []ProductVM {
	vms := make([]ProductVM, len(products))
	for i, p := range products {
		vms[i] = NewProductVM(p)
	}
	return vms
}

// NewTransactionVM maps a models.PosTransaction to a TransactionVM.
func NewTransactionVM(tx models.PosTransaction) TransactionVM {
	return TransactionVM{
		IDShort:              tx.ID[:12] + "...",
		PaymentMethodLabel:   string(tx.PaymentMethod),
		IsCash:               tx.PaymentMethod == models.PaymentMethodCash,
		TotalAmountFormatted: "Rp " + utils.FormatCurrency(tx.TotalAmount),
		CreatedAtFormatted:   tx.CreatedAt.Format("02 Jan 2006, 15:04") + " WIB",
	}
}

// NewTransactionVMs maps a slice of models.PosTransaction to a slice of TransactionVM.
func NewTransactionVMs(transactions []models.PosTransaction) []TransactionVM {
	vms := make([]TransactionVM, len(transactions))
	for i, tx := range transactions {
		vms[i] = NewTransactionVM(tx)
	}
	return vms
}

// stockBadge returns the display label and Tailwind CSS class for a given stock level.
func stockBadge(stock int) (label, class string) {
	switch {
	case stock > 10:
		return utils.IntToString(stock) + " Available", "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
	case stock > 0:
		return utils.IntToString(stock) + " Low Stock", "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-800"
	default:
		return "Out of Stock", "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800"
	}
}

// FormatTime formats a time.Time value for display in Indonesian locale.
// This is a standalone helper for use in places where full VM mapping is not needed.
func FormatTime(t time.Time) string {
	return t.Format("02 Jan 2006, 15:04") + " WIB"
}
