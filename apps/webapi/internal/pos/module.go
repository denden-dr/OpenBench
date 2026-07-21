package pos

// Module encapsulates the external-facing handlers of the pos package.
type Module struct {
	Handler *Handler
	Service Service
}

// NewModule initializes the pos module with a ready service.
func NewModule(svc Service) Module {
	return Module{
		Handler: NewHandler(svc),
		Service: svc,
	}
}
