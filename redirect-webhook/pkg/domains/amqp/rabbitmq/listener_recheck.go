package rabbitmq

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
)

type RecheckConnectionsListener struct {
	handleMessage func(tenantKey *tenants.TenantKey)
}

func (c RecheckConnectionsListener) ConsumeMessages() {

	logging.Logger.Info("[RABBIT] Listening to RecheckConnections queue")

	deliveries, err := manager.channel.Consume(
		formatQueueName(QueueRecheckConnections),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register RecheckConnections queue")

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
					splunk.ConnectionManagementLog{}.SendError("error handling connection recheck message", e, nil)
				}
			}
		}(deliveries)
	}

	logging.Logger.Info("[RABBIT] Waiting for messages in the queue: " + QueueRecheckConnections)
}
