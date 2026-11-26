package rabbitmq

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
)

type FlushConnectionsListener struct {
	handleMessage func(tenantKey *tenants.TenantKey)
}

func (c FlushConnectionsListener) ConsumeMessages() {

	logging.Logger.Info("[RABBIT] Listening to FlushConnections queue")

	deliveries, err := manager.channel.Consume(
		formatQueueName(QueueConnectionsFlush),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register FlushConnections queue")

	if err == nil {
		go func(deliveries <-chan amqp.Delivery) {

			for message := range deliveries {
				e := func(delivery amqp.Delivery) error {
					tenantKey := &tenants.TenantKey{}
					bytes := delivery.Body
					e := json.Unmarshal(bytes, tenantKey)
					if e != nil {
						return errors.New("error unmarshalling message on connection management queue: " + e.Error())
					}

					c.handleMessage(tenantKey)

					return nil
				}(message)
				if e != nil {
					splunk.ConnectionManagementLog{}.SendError("error handling connection flush message", e, nil)
				}
			}
		}(deliveries)
	}

	logging.Logger.Info("[RABBIT] Waiting for messages in the queue: " + QueueConnectionsFlush)
}
