package delivery_queues

import (
	. "conductor.com.br/pier/hooks/pkg/error"
)

type Error struct {
	OrmHandler
}

func (e Error) HandleEventTypeAlreadyConfigured() {
	HttpError{Status: 400, Description: "event type is already configured"}.Throw()
}