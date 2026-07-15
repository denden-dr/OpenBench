package models

import "time"

type PaymentMethod string

const (
	PaymentMethodCash PaymentMethod = "CASH"
	PaymentMethodQRIS PaymentMethod = "QRIS"
)

type PosTransaction struct {
	ID            string               `json:"id" db:"id"`
	PaymentMethod PaymentMethod        `json:"payment_method" db:"payment_method"`
	TotalAmount   int64                `json:"total_amount" db:"total_amount"`
	CreatedAt     time.Time            `json:"created_at" db:"created_at"`
	Items         []PosTransactionItem `json:"items,omitempty"`
}

type PosTransactionItem struct {
	ID            string `json:"id" db:"id"`
	TransactionID string `json:"transaction_id" db:"transaction_id"`
	ProductID     string `json:"product_id" db:"product_id"`
	Quantity      int    `json:"quantity" db:"quantity"`
	Price         int64  `json:"price" db:"price"`
	ProductName   string `json:"product_name,omitempty" db:"product_name"` // populated dynamically on queries
}

type CheckoutItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type CheckoutRequest struct {
	PaymentMethod PaymentMethod         `json:"payment_method" validate:"required,oneof=CASH QRIS"`
	Items         []CheckoutItemRequest `json:"items" validate:"required,min=1,dive"`
}
