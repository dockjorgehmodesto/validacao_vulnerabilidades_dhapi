package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"github.com/go-chi/chi"
)

type Internals struct {
	Tenants helpers.Router
}

func (i Internals) Route(r chi.Router) {
	r.Route("/tenants", i.Tenants.Route)
}
