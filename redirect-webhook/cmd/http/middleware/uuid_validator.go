package middleware

import (
	"net/http"

	httphelpers "conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/go-chi/chi"
)

func ValidateUUID(wildcard string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, wildcard)
			if valid, err := helpers.ValidatorUUID(id); err != nil || !valid {
				httphelpers.HandleInvalidUUID(err)
			}
			next.ServeHTTP(w, r)
		})
	}
}
