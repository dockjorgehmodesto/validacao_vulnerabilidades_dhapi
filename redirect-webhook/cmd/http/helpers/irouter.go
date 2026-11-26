package helpers

import "github.com/go-chi/chi"

type Router interface {
	Route(r chi.Router)
}
