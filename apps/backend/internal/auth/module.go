package auth

import (
	"time"

	"github.com/denden-dr/OpenBench/apps/backend/config"
	"github.com/jmoiron/sqlx"
)

// Module encapsulates the external-facing handlers and query repositories of the auth package.
type Module struct {
	Handler   *Handler
	QueryRepo QueryRepository
}

// NewModule initializes the entire auth domain, spins up the cleanup worker,
// and returns the Module along with a cleanup stop function.
func NewModule(db *sqlx.DB, cfg *config.Config) (Module, func()) {
	qr := NewQueryRepository(db)
	cr := NewCommandRepository(db)
	svc := NewService(qr, cr, cfg)

	w := NewCleanupWorker(cr, 24*time.Hour)
	w.Start()

	return Module{
		Handler:   NewHandler(svc, cfg),
		QueryRepo: qr,
	}, w.Stop
}
