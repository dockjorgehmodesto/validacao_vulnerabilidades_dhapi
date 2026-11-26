package connectors

import (
	"context"
	servicebus "github.com/Azure/azure-service-bus-go"
)

// Interface with common methods used by Azure's ServiceBus
type ServiceBusSubscription interface {
	servicebus.Closer
	Receive(ctx context.Context, handler servicebus.Handler) error
	NewDeadLetter() *servicebus.DeadLetter
	PeekOne(ctx context.Context, options ...servicebus.PeekOption) (*servicebus.Message, error)
}
