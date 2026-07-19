package warranty

import (
	"github.com/denden-dr/OpenBench/internal/database"
	"github.com/jmoiron/sqlx"
)

// Module encapsulates the external-facing handlers and services of the warranty package.
type Module struct {
	Handler    *Handler
	Service    Service
}

// NewModule initializes the entire warranty domain layer.
func NewModule(db *sqlx.DB, txManager database.TxManager) Module {
	qr := NewQueryRepository(db)
	cr := NewCommandRepository(db)
	svc := NewService(qr, cr, txManager)
	return Module{
		Handler:    NewHandler(svc),
		Service:    svc,
	}
}
