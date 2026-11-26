package rabbitmq

import (
	"context"
	"errors"
	"strconv"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"github.com/streadway/amqp"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	RabbitMQContentType    = "text/plain"
	RabbitMQExchangeTopic  = "topic"
	RabbitMQExchangeFanout = "fanout"
)

// Publish : Message delivery method on a specific defaultChannel
func PublishString(message string, exchange string, routingKey string) {

	maxAttempts := global.Viper.GetInt("rabbit.headers.attempts")

	go publishAsync(message, exchange, routingKey, maxAttempts, 1)
}

// PublishAsyncWithAttempts :
func PublishWithAttempts(message string, exchange string, routingKey string, maxAttempts int, attempts int) {
	go publishAsync(message, exchange, routingKey, maxAttempts, attempts)
}

func publishAsync(message string, exchange string, routingKey string, maxAttempts int, attempts int) {

	Publish([]byte(message), exchange, routingKey, maxAttempts, attempts)
}

func Publish(message []byte, exchange string, routingKey string, maxAttempts int, attempts int) {
	go func() {
		headers := amqp.Table{
			"x-max-attempts": strconv.Itoa(maxAttempts),
			"attempts":       attempts,
		}

		publishing := amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  RabbitMQContentType,
			Body:         message,
			Headers:      headers,
		}

		if manager.channel == nil {
			logging.Logger.Error("[RABBIT] channel not available")
			return
		}

		err := manager.channel.Publish(
			exchange,
			routingKey,
			false,
			false,
			publishing)

		failOnError(err, "Failed to publish a message")
	}()
}

func PublishSync(message []byte, exchange string, routingKey string, maxAttempts int, attempts int, ctx context.Context) error {
	ddTracer, ctx := tracer.StartSpanFromContext(ctx, "PublishSyncMessageSentToRabbit")
	defer ddTracer.Finish()

	headers := amqp.Table{
		"x-max-attempts": strconv.Itoa(maxAttempts),
		"attempts":       attempts,
	}

	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  RabbitMQContentType,
		Body:         message,
		Headers:      headers,
	}

	if manager.channel == nil {
		return errors.New("channel not available")
	}

	return manager.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		publishing)
}

func delay(attempts int) int {
	if attempts > 5 {
		return milliSecond(5)
	}
	return milliSecond(attempts)
}

func milliSecond(value int) int {
	return value * 60 * 1000
}
