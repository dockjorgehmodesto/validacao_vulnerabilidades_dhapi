package rabbitmq

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"github.com/streadway/amqp"
)

// Exchange names
const (
	ExchangeRecheckConnections = "CONNECTION_RECHECK_FANOUT"
	ExchangeConnectionFlush    = "CONNECTION_FLUSH_FANOUT"
)

// Queue
const (
	QueueRecheckConnections    = "pier.hooks.CONNECTIONS_RECHECK"
	QueueRecheckConnectionsKey = "pier.hooks.CONNECTIONS_RECHECK_KEY"
	QueueConnectionsFlush      = "pier.hooks.CONNECTIONS_FLUSH"
	QueueConnectionsFlushKey   = "pier.hooks.CONNECTIONS_FLUSH_KEY"
)

type ListenerActions struct {
	RecheckConnections func(tenantKey *tenants.TenantKey)
	FlushConnections   func(tenantKey *tenants.TenantKey)
}

// Consumers : List of message consumers
func (l ListenerActions) Build() []Listener {
	return []Listener{
		RecheckConnectionsListener{l.RecheckConnections},
		FlushConnectionsListener{l.FlushConnections},
	}
}

// Listener : RabbitAMQ message queue consumer method
type Listener interface {
	ConsumeMessages()
}

// Exchange : Representation structure of a message exchange defaultChannel
type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// Queue : Representation of the structure of a message queue
type Queue struct {
	Name        string
	Key         string
	Exchange    string
	Durable     bool
	AutoDelete  bool
	Exclusive   bool
	NoWait      bool
	Args        amqp.Table
	IsBroadcast bool
}

// Exchanges : List of message defaultChannel settings
func Exchanges() []Exchange {
	return []Exchange{
		{ExchangeRecheckConnections, RabbitMQExchangeFanout, true, false, false, false, nil},
		{ExchangeConnectionFlush, RabbitMQExchangeFanout, true, false, false, false, nil},
	}
}

// Queues : List of Listener Queue Settings
func Queues() []Queue {
	args := amqp.Table{"x-message-ttl": 600000}
	return []Queue{
		{QueueRecheckConnections, QueueRecheckConnectionsKey, ExchangeRecheckConnections, false, true, true, false, args, true},
		{QueueConnectionsFlush, QueueConnectionsFlushKey, ExchangeConnectionFlush, false, true, true, false, args, true},
	}
}
