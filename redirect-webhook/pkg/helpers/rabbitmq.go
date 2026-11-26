package helpers

import (
	"context"
	"github.com/streadway/amqp"
	"time"
)

func ReadRabbitMsgsInGoroutines(maxGoroutines int, readMessageFunc func(ctx context.Context, cancelFunc context.CancelFunc, message amqp.Delivery), deliveries <-chan amqp.Delivery) {
	go func(deliveries <-chan amqp.Delivery) {
		guard := make(chan struct{}, maxGoroutines)

		for message := range deliveries {

			// RabbitMQ can lock consumers due to internal issues (e.g. memory limit)
			// This can cause locks in the "readMessageFunc", to prevent it, we will
			// use a context with a timeout, guaranteeing the "guard" channel is cleared.
			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)

			guard <- struct{}{} // Blocks when channel is full
			go func(ctx context.Context, cancelFunc context.CancelFunc, message amqp.Delivery) {
				readMessageFunc(ctx, cancelFunc, message)

				// Clearing the guard channel, allowing another goroutine to start
				// This will happen when the cancelFunc is called or when there's a timeout in the context
				select {
				case <-ctx.Done():
					<-guard // Releases goroutine for another execution
				}
			}(ctx, cancelFunc, message)

		}
	}(deliveries)
}
