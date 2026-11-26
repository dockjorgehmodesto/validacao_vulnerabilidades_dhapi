package message_processors

import "errors"

var (
	WebHookNotSetup    = errors.New("event with no webhook setup")
	EncryptionNotFound = errors.New("attempt to trigger webhook without encryption")
)
