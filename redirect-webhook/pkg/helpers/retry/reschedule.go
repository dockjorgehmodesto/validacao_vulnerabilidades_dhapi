package retry

import (
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"errors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"math"
	"time"
)

var RetryLimit = viper.GetInt32("rabbit.retryNumber")

type DelayFunc func(retryCounter int32) (delayMillis int64)

type RescheduleOptions struct {
	Channel      *amqp.Channel
	ExchangeName string
	RoutingKey   string
}

/*
Increments delays exponentially with a PG
*/
func DefaultDelayFunc(retryCounter int32) int64 {
	return int64((1 * math.Pow(3.21, float64(retryCounter-1))) * 1000)
}

/*
Auxiliary method for rescheduling using the DefaultDelayFunc
*/
func RescheduleMessage(delivery amqp.Delivery, opts *RescheduleOptions) error {
	return RescheduleMessageWithDelayFunc(DefaultDelayFunc, delivery, opts)
}

/*
Behaves the same as RescheduleMessageDeclaringDelay, working with a DelayFunc
*/
func RescheduleMessageWithDelayFunc(delayFunc DelayFunc, delivery amqp.Delivery, opts *RescheduleOptions) error {
	var retries int32 = 1
	if delivery.Headers != nil {
		if mapVal, ok := delivery.Headers["retries"]; ok {
			retries = mapVal.(int32)
			retries = retries + 1
		}
	}

	if retries > RetryLimit {
		logging.Logger.Debug("retry limit", zap.Int32("retryCount", retries))
		return nil
	}

	delayMillis := delayFunc(retries)
	logging.Logger.Debug("retry number", zap.Int32("retryCount", retries), zap.Int64("delayMillis", delayMillis))

	delivery.Headers = setupHeaders(delayMillis, delivery.Headers, retries)
	return RescheduleMessageDelay(delayMillis, delivery, opts)
}

/*
Copies a message body and headers, updating the headers used by the library: "x-delay", "retries" and "message-uuid".
Once setup, the message is encapsulated in an amqp.Publishing and resent to the exchange.
A message will be available again in the queue after the given delayMillis
Delivery represents the message that should be rescheduled
Channels must be active for rescheduling
The exchangeName must follow the one given on the CreateDefaultDelayExchange method
The routingKey is present inside the amqp.Delivery structure, still it does not always seem to match. This will also
allow custom routing
*/
func RescheduleMessageDelay(delayMillis int64, delivery amqp.Delivery, opts *RescheduleOptions) error {
	if opts == nil {
		return errors.New("reschedule options are invalid")
	}

	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  delivery.ContentType,
		Body:         delivery.Body,
		Headers:      delivery.Headers,
	}

	if opts.Channel == nil {
		return errors.New("RabbitMQ not available")
	}

	// Both flags bellow set as false helps guaranteeing the message will be delivered
	return opts.Channel.Publish(
		opts.ExchangeName,
		opts.RoutingKey,
		false,
		false,
		publishing,
	)
}

func setupHeaders(delayMillis int64, deliveryHeaders amqp.Table, retries int32) amqp.Table {
	var headers amqp.Table
	if deliveryHeaders != nil && len(deliveryHeaders) != 0 {
		headers = deliveryHeaders
	} else {
		headers = make(amqp.Table)
	}

	headers["x-delay"] = delayMillis
	headers["retries"] = retries

	if _, ok := headers["message-uuid"]; !ok {
		headers["message-uuid"] = uuid.New().String()
	}
	return headers
}
