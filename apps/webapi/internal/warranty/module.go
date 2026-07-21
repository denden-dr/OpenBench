package warranty

// Module encapsulates the external-facing handlers and services of the warranty package.
type Module struct {
	Handler *Handler
	Service Service
}

// NewModule initializes the warranty module with a ready service.
func NewModule(svc Service) Module {
	return Module{
		Handler: NewHandler(svc),
		Service: svc,
	}
}
