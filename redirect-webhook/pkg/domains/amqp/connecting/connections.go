package connecting

import (
	"conductor.com.br/pier/hooks/pkg/helpers"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	AMQPBrokerServiceBus = "SERVICEBUS"
	AMQPBrokerRabbitMQ   = "RABBITMQ"
)

type Connection struct {
	helpers.GormHooks
	ID               string
	ConnectionString string
	Enabled          bool
	IssuerId         int
	Broker           string
	EventType        string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Properties for connecting to ServiceBus
	Topic            string
	SubscriptionName string

	// Properties for connecting to RabbitMQ
	VirtualHost   string
	ExchangeName  string
	QueueName     string
	QueueKey      string
	ProcessorType string
}

func (c Connection) TableName() string {
	return "amqp_connections"
}

func (c Connection) ValidatePropertiesForBroker() error {
	// All kinds of brokers will require a connectionString
	if c.ConnectionString == "" {
		return errors.New("connectionString/dial not informed")
	}

	switch c.Broker {
	case AMQPBrokerServiceBus:
		if c.Topic == "" {
			return errors.New("topic not informed")
		}
		if c.SubscriptionName == "" {
			return errors.New("subscriptionName not informed")
		}
		break
	case AMQPBrokerRabbitMQ:
		if c.VirtualHost == "" {
			return errors.New("virtualHost not informed")
		}
		if c.ExchangeName == "" {
			return errors.New("exchangeName not informed")
		}
		if c.QueueName == "" {
			return errors.New("queueName not informed")
		}
		if c.QueueKey == "" {
			return errors.New("queueKey not informed")
		}
		break
	default:
		return errors.New("invalid broker")
	}
	return nil
}

// Omits sensitive data in a connection string
func (c Connection) GetMaskedConnectionString() string {
	switch c.Broker {
	case AMQPBrokerServiceBus:
		return c.maskConnectionStringForServiceBus()
	case AMQPBrokerRabbitMQ:
		return c.maskConnectionStringForRabbitMQ()
	default:
		return ""
	}
}

func (c Connection) maskConnectionStringForServiceBus() string {
	result := make(map[string]string)
	stringKeyMapSlice := strings.Split(c.ConnectionString, ";")
	for _, stringKeyMap := range stringKeyMapSlice {
		pairString := strings.Split(stringKeyMap, "=")
		if len(pairString) == 2 {
			result[pairString[0]] = pairString[1]
		}
	}

	return fmt.Sprintf("Endpoint=%s;SharedAccessKeyName=%s;SharedAccessKey=********", result["Endpoint"], result["SharedAccessKeyName"])
}

func (c Connection) maskConnectionStringForRabbitMQ() string {
	splitAt := strings.Split(c.ConnectionString, "@")
	removeAmqp := strings.Split(splitAt[0], "amqp://")[1]
	credentials := strings.Split(removeAmqp, ":")
	user := credentials[0]

	return fmt.Sprintf("amqp://%s:*****@%s", user, splitAt[1])
}
