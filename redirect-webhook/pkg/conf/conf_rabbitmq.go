package conf

import (
	"conductor.com.br/pier/hooks/pkg/domains"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/rabbitmq"
)

func ConfigRabbitMQ(services domains.ServiceContainer) {
	listeners := rabbitmq.ListenerActions{
		RecheckConnections: services.AMQP.RecheckConnections,
		FlushConnections:   services.AMQP.ReloadAndFlush,
	}
	listeners.Build()
	go rabbitmq.InitConnections(listeners)
}
