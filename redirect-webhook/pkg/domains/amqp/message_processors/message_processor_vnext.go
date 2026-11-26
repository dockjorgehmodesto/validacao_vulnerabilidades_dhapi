package message_processors

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

type MessageProcessorVNext struct {
	DefaultMessageProcessor
}

type VNextMessage struct {
	EventType string    `json:"eventType"`
	Timestamp time.Time `json:"timestamp"`
	IssuerID  int       `json:"issuerId"`
}

func (vNext MessageProcessorVNext) Process(ctx context.Context, message Message) (int, error) {

	vNextMessage := VNextMessage{}
	err := json.Unmarshal([]byte(message.Data), &vNextMessage)
	if err != nil {
		return 0, errors.Wrap(err, "could not unmarshal")
	}
	message.EventType = vNextMessage.EventType

	return vNext.DefaultMessageProcessor.process(ctx, message, vNextMessage.IssuerID)
}
