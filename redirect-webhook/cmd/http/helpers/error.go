package helpers

import . "conductor.com.br/pier/hooks/pkg/error"

type HandlerError struct {
}

func (h HandlerError) JsonDecode() *HttpError {
	return &HttpError{Status: 400, Description: "error decoding json", DisplayStack: true}
}

func (h HandlerError) JsonMarshal() *HttpError {
	return &HttpError{Status: 400, Description: "error generating json", DisplayStack: true}
}

func (h HandlerError) JsonValidation() *HttpError {
	return &HttpError{Status: 400, Description: "error validating json", DisplayStack: true}
}

func (h HandlerError) InvalidUUID() *HttpError {
	return &HttpError{Status: 400, Description: "informed id is not a valid UUID", DisplayStack: true}
}
