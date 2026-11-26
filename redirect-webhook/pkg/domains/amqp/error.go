package amqp

import (
	. "conductor.com.br/pier/hooks/pkg/error"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"github.com/jinzhu/gorm"
)

type Error struct {
	OrmHandler
}

func (e Error) HandleDisconnectionError(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "error disconnecting to AMQP connection"}.WithChild(err).Throw()
	}
}

func (e Error) HandleConnectionFlushError(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "error flushing AMQP connection"}.WithChild(err).Throw()
	}
}

// Should raise a panic if an topic is found(if there's no error), or if it happened a DB connection error
func (e Error) HandleEventTypeValidationOnCreate(err error) {
	if err == nil {
		HttpError{Status: 400, Description: "AMQP connection with given eventType already exists"}.Throw()
	}

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		HttpError{Status: 500, Description: "error validating topic"}.WithChild(err).Throw()
	}
}

func (e Error) HandleListEnabledConnectionsOnRecheck(err error) {
	splunk.ConnectionManagementLog{}.SendError("could not list enabled connections on recheck", err, nil)
}

func (e Error) HandleListEnabledConnectionsOnReset(err error) {
	splunk.ConnectionManagementLog{}.SendError("could not list enabled connections when resetting", err, nil)
}

func (e Error) HandleDisconnectOnRecheck(err error) {
	splunk.ConnectionManagementLog{}.SendError("could not disconnect to a disabled connection", err, nil)
}

func (e Error) HandleInvalidEventType(isValid bool) {
	if !isValid {
		HttpError{Status: 400, Description: "event type is invalid"}.Throw()
	}
}

func (e Error) HandleConnectionTestError(err error) {
	if err != nil {
		HttpError{Status: 400, Description: "could not connect to message broker"}.WithChild(err).Throw()
	}
}

func (e Error) HandleActiveConnectionNotFoundByEventType(err error) {
	if err != nil {
		HttpError{Status: 400, Description: "informed event type is not associated to an active connection"}.WithChild(err).Throw()
	}
}

func (e Error) HandleResend(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "could not resend messages", DisplayStack: true}.WithChild(err).Throw()
	}
}

func (e Error) HandleEventWithoutEncryptionKey(err error) {
	if err != nil {
		HttpError{Status: 400, Description: "issuer not setup for encryption"}.WithChild(err).Throw()
	}
}

func (e Error) HandleEventWithoutWebHook(err error) {
	if err != nil {
		HttpError{Status: 400, Description: "event not associated to a webhook"}.WithChild(err).Throw()
	}
}

func (e Error) HandleMessageCount(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "could not count how many messages are available"}.WithChild(err).Throw()
	}
}

func (e Error) HandleInvalidMessageAmount(msgCount int) {
	if msgCount <= 0 {
		HttpError{Status: 400, Description: "no messages to send"}.Throw()
	}
}

func (e Error) HandleOnFetchMetadata(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "could not fetch metadata for the connection"}.WithChild(err).Throw()
	}
}

func (e Error) HandleInvalidBrokerProperties(err error) {
	if err != nil {
		HttpError{Status: 400, Description: err.Error()}.WithChild(err).Throw()
	}
}

func (e Error) HandleInvalidMessageProcessor(err error) {
	HttpError{Status: 400, Description: "invalid message processor"}.WithChild(err).Throw()
}

func (e Error) HandleMapUpdate(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "could not parse update map"}.WithChild(err).Throw()
	}
}
