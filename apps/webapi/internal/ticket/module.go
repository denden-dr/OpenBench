package ticket

// Module encapsulates the external-facing handlers of the ticket package.
type Module struct {
	Handler *Handler
	Service Service
}

// NewModule initializes the ticket module with a ready service.
func NewModule(svc Service) Module {
	return Module{
		Handler: NewHandler(svc),
		Service: svc,
	}
}
