package inventory

// Module encapsulates the external-facing handlers and repositories of the inventory package.
type Module struct {
	Handler     *Handler
	Service     Service
	QueryRepo   QueryRepository
	CommandRepo CommandRepository
}

// NewModule initializes the inventory module with ready service and repos.
func NewModule(svc Service, qr QueryRepository, cr CommandRepository) Module {
	return Module{
		Handler:     NewHandler(svc),
		Service:     svc,
		QueryRepo:   qr,
		CommandRepo: cr,
	}
}
