package auth

import (
	"time"

	"github.com/denden-dr/OpenBench/config"
	"github.com/jmoiron/sqlx"
	"github.com/samber/hot"
)

// Module encapsulates the external-facing handlers and query repositories of the auth package.
type Module struct {
	Handler    *Handler
	WebHandler *WebHandler
	QueryRepo  QueryRepository
}

// NewModule initializes the entire auth domain, spins up the cleanup worker,
// and returns the Module along with a cleanup stop function.
func NewModule(db *sqlx.DB, cfg *config.Config) (Module, func()) {
	cache := hot.NewHotCache[string, bool](hot.WTinyLFU, 10000).
		WithTTL(15 * time.Minute).
		WithJanitor().
		Build()

	qr := NewQueryRepository(db, cache)
	cr := NewCommandRepository(db, cache)
	svc := NewService(qr, cr, cfg)

	w := NewCleanupWorker(cr, 24*time.Hour)
	w.Start()

	cleanup := func() {
		w.Stop()
		cache.StopJanitor()
	}

	return Module{
		Handler:    NewHandler(svc, cfg),
		WebHandler: NewWebHandler(svc, cfg),
		QueryRepo:  qr,
	}, cleanup
}
