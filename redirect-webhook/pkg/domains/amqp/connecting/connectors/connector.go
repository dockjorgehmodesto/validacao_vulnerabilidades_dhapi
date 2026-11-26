package connectors

import (
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/message_processors"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"errors"
	"sync"
	"time"
)

type Connector interface {
	// Tests whether a connection data is valid
	Test() error

	Listen()

	Disconnect() error

	GetEventType() string

	GetId() string

	GetIssuerId() int

	SetOnClose(onclose func(key string))

	GetConnection() *connecting.Connection

	GetProcessorType() string
}

func CreateConnector(connection connecting.Connection, tenantKey *tenants.TenantKey, processor message_processors.MessageProcessor) (Connector, error) {
	switch connection.Broker {
	case connecting.AMQPBrokerRabbitMQ:
		return &rabbitMQConnector{
			Connection:      connection,
			reconnectMutex:  sync.Mutex{},
			ConnectedAt:     time.Now(),
			tenantKey:       tenantKey,
			processor:       processor,
			isReconnecting:  false,
			isDisconnecting: false,
			disconnectMutex: sync.Mutex{},
		}, nil
	case connecting.AMQPBrokerServiceBus:
		return &serviceBusConnector{
			Connection:      connection,
			mutex:           sync.Mutex{},
			ConnectedAt:     time.Now(),
			tenantKey:       tenantKey,
			processor:       processor,
			forceClose:      false,
			isDisconnecting: false,
		}, nil

	default:
		return nil, errors.New("invalid broker")
	}
}

func shouldRetryMessage(e error) bool {
	return e != nil && !(e == message_processors.WebHookNotSetup || e == message_processors.EncryptionNotFound)
}
