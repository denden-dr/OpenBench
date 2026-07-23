package dashboard

import "github.com/jmoiron/sqlx"

type Module struct {
	Handler *Handler
	Service Service
	Repo    Repository
}

func NewModule(db *sqlx.DB) Module {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	return Module{
		Handler: handler,
		Service: svc,
		Repo:    repo,
	}
}
