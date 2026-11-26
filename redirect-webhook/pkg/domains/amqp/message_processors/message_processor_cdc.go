package message_processors

import (
	"context"
)

type MessageProcessorCDC struct {
	DefaultMessageProcessor
}

func (cdc MessageProcessorCDC) Process(ctx context.Context, message Message) (int, error) {
	return cdc.DefaultMessageProcessor.process(ctx, message, cdc.connection.IssuerId)
}
