package middleware

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers/constants"
	"context"
	"net/http"
	"strconv"
)

func CheckIssuerId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		issuerId, _ := strconv.Atoi(r.Header.Get(constants.ConstHeaderIssuerId))

		if issuerId <= 0 {
			err := MiddlewareError{}.IssuerIdNotFound()
			helpers.JsonResponse(w, err, err.Status)
			return
		} else {

			tenantKey, _, e := tenants.TenantService.FetchOnMapByIssuerId(issuerId)
			if e != nil {
				err := MiddlewareError{}.IssuerDoesNotHavePermission()
				helpers.JsonResponse(w, err, err.Status)
				return
			} else {
				ctx := context.WithValue(r.Context(), constants.ConstHeaderIssuerId, issuerId)
				ctx = context.WithValue(ctx, constants.ConstTenantHeaderContext, tenantKey)
				next.ServeHTTP(w, r.WithContext(ctx))
			}

		}
	})
}
