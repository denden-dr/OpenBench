package pos

import (
	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/jmoiron/sqlx"
)

// Module encapsulates the external-facing handlers of the pos package.
type Module struct {
	Handler    *Handler
	WebHandler *WebHandler
}

// NewModule initializes the entire pos domain layer.
func NewModule(
	db *sqlx.DB,
	txManager database.TxManager,
	invReader InventoryProductReader,
	invWriter InventoryStockWriter,
) Module {
	qr := NewQueryRepository(db)
	cr := NewCommandRepository(db)
	svc := NewService(qr, cr, invReader, invWriter, txManager)
	return Module{
		Handler:    NewHandler(svc),
		WebHandler: NewWebHandler(svc, invReader),
	}
}
