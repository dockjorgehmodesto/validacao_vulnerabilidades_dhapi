package error

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

type UnhandledError struct {
	Description string `json:"description"`
}

func (u UnhandledError) Error() string {
	return u.Description
}

type HttpError struct {
	Status          int     `json:"-"`
	Description     string  `json:"description,omitempty"`
	Children        []error `json:"children,omitempty"`
	DisplayStack    bool    `json:"-"`
	IsDatabaseError *bool   `json:"isDatabaseError,omitempty"`
}

func (h HttpError) ShouldDisplayStack() bool {
	if global.Viper.GetString("app.profile") == "development" {
		return true
	}

	return h.DisplayStack
}

func (h HttpError) Throw() {
	log.Error(h)
	panic(h)
}

func (h HttpError) ThrowWithMessage(msg ...string) {
	h.Description = fmt.Sprintf(h.Description, msg)
	h.Throw()
}

// PrintStackTrace : Recovers StackTrace from Error Thrown
func (h HttpError) GetStackTrace() string {
	return fmt.Sprintf("%+v", h.Children)
}

// Error : retrieve error message thrown
func (h HttpError) Error() string {
	return h.Description
}

func (h HttpError) WithChild(e error) *HttpError {
	if e != nil {
		if len(h.Children) == 0 {
			h.Children = make([]error, 0)
		}
		h.Children = append(h.Children, UnhandledError{Description: e.Error()})
	}
	return &h
}

func (h HttpError) HandleInternalServerError(e error) {
	if e != nil {
		h.Status = http.StatusInternalServerError
		h.Description = http.StatusText(h.Status)
		b := false
		h.IsDatabaseError = &b
		h.WithChild(e).Throw()
	}
}

func (h HttpError) HandleDatabaseError(e error) {
	if e != nil {
		h.Status = http.StatusInternalServerError
		h.Description = http.StatusText(h.Status)
		b := true
		h.IsDatabaseError = &b
		h.WithChild(e).Throw()
	}
}

func (h HttpError) HandleDatabaseErrorWithMessage(e error, message string) {
	if e != nil {
		h.Status = http.StatusInternalServerError
		h.Description = message
		b := true
		h.IsDatabaseError = &b
		h.WithChild(e).Throw()
	}
}

func (h HttpError) BadRequest(msg string) HttpError {
	return HttpError{Status: http.StatusBadRequest, Description: msg}
}
