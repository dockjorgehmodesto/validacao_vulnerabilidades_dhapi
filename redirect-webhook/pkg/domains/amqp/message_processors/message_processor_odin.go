package message_processors

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

type MessageProcessorOdin struct {
	DefaultMessageProcessor
}

type OdinMessage struct {
	MessageId string `json:"uuidTransacao"`
	IssuerId  int    `json:"idEmissor,string"`
}

func (odin MessageProcessorOdin) Process(ctx context.Context, message Message) (int, error) {

	payload := OdinMessage{}
	err := json.Unmarshal([]byte(message.Data), &payload)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	message.Id = payload.MessageId
	return odin.DefaultMessageProcessor.process(ctx, message, payload.IssuerId)
}
