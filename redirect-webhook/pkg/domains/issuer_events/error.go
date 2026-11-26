package issuer_events

import (
	. "conductor.com.br/pier/hooks/pkg/error"
)

type Error struct {
	OrmHandler
}

func (e Error) HandleNameAlreadyTaken() {
	HttpError{Status: 400, Description: "event name is already taken"}.Throw()
}

func (e Error) HandleInvalidConnection() {
	HttpError{Status: 400, Description: "connection id is invalid"}.Throw()
}

func (e Error) HandleConnectionAlreadyAssociatedToAnEvent() {
	HttpError{Status: 400, Description: "connection id is already associated to an event for this issuer"}.Throw()
}

func (e Error) HandleEventNotSetupForIssuer() {
	HttpError{Status: 400, Description: "event not available for issuer"}.Throw()
}
