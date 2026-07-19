package ticket

import (
	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/denden-dr/OpenBench/internal/events"
	"github.com/jmoiron/sqlx"
)

// Module encapsulates the external-facing handlers of the ticket package.
type Module struct {
	Handler    *Handler
}

// NewModule initializes the entire ticket domain layer.
func NewModule(
	db *sqlx.DB,
	txManager database.TxManager,
	warrantyGen WarrantyGenerator,
	eventBus events.EventBus,
	encryptionKey string,
) Module {
	qr := NewQueryRepository(db)
	cr := NewCommandRepository(db)
	svc := NewService(qr, cr, txManager, warrantyGen, eventBus, encryptionKey)
	return Module{
		Handler:    NewHandler(svc),
	}
}
