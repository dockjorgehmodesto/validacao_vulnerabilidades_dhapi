package helpers

import (
	"encoding/json"
	"net/http"
)

func HandleInvalidUUID(e error) {
	HandlerError{}.InvalidUUID().WithChild(e).Throw()
}

func addDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func Ok(w http.ResponseWriter, data interface{}) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusOK)
	if bytes, e := json.Marshal(data); e == nil {
		if _, e := w.Write(bytes); e != nil {
			HandlerError{}.JsonDecode().Throw()
		}
	} else {
		HandlerError{}.JsonDecode().Throw()
	}
}

func Success(w http.ResponseWriter) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func JsonResponse(w http.ResponseWriter, data interface{}, httpStatus int) {
	addDefaultHeaders(w)
	w.WriteHeader(httpStatus)
	if bytes, e := json.Marshal(data); e != nil {
		HandlerError{}.JsonDecode().Throw()
	} else {
		if _, e := w.Write(bytes); e != nil {
			HandlerError{}.JsonDecode().Throw()
		}
	}
}

func NoContent(w http.ResponseWriter) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}
