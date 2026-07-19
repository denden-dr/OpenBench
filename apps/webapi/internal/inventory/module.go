package inventory

import "github.com/jmoiron/sqlx"

// Module encapsulates the external-facing handlers and repositories of the inventory package.
type Module struct {
	Handler     *Handler
	QueryRepo   QueryRepository
	CommandRepo CommandRepository
}

// NewModule initializes the entire inventory domain layer.
func NewModule(db *sqlx.DB) Module {
	qr := NewQueryRepository(db)
	cr := NewCommandRepository(db)
	svc := NewService(qr, cr)
	return Module{
		Handler:     NewHandler(svc),
		QueryRepo:   qr,
		CommandRepo: cr,
	}
}
