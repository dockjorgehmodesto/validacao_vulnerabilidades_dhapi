package middleware

import . "conductor.com.br/pier/hooks/pkg/error"

type MiddlewareError struct {
}

func (m MiddlewareError) IssuerIdNotFound() *HttpError {
	return &HttpError{Status: 403, Description: "issuer ID not found on request headers", DisplayStack: false}
}

func (m MiddlewareError) IssuerDoesNotHavePermission() *HttpError {
	return &HttpError{Status: 403, Description: "issuer doesn't have permission for this feature", DisplayStack: false}
}
