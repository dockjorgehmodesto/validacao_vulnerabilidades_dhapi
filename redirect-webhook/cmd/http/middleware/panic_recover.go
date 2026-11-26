package middleware

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	errors "conductor.com.br/pier/hooks/pkg/error"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				var response errors.HttpError
				if value, ok := rvr.(errors.HttpError); ok {
					response = value
				} else if value, ok := rvr.(*errors.HttpError); ok {
					response = *value
				} else {

					var e error
					if errorValue, ok := rvr.(error); ok {
						e = errorValue
					} else {
						e = errors.UnhandledError{Description: fmt.Sprintf("%v", rvr)}
					}
					response = errors.HttpError{
						Status:       500,
						Description:  "Internal server error",
						DisplayStack: false,
					}
					response = *response.WithChild(e)
				}

				splunkCtxPointer := r.Context().Value("splunkContext")
				if splunkCtxPointer != nil {
					if splunkCtx, ok := splunkCtxPointer.(*splunk.SplunkContext); ok {
						splunkCtx.Stack = rvr
						ctx := context.WithValue(r.Context(), "splunkContext", &splunkCtx)
						r.WithContext(ctx)
					}
				} else {
					log.Errorf("nil splunkContext on serving route [%s] %s", r.Method, r.RequestURI)
				}

				if !response.ShouldDisplayStack() {
					response.Children = nil
					response.IsDatabaseError = nil
				}

				if response.Status == http.StatusNoContent {
					helpers.NoContent(w)
				} else {
					helpers.JsonResponse(w, response, response.Status)
				}
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
