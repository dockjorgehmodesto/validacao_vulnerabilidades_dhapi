package message_processors

import (
	"errors"

	"conductor.com.br/pier/hooks/pkg/domains/abstraction"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/delivery_queues"
	"conductor.com.br/pier/hooks/pkg/domains/encryption"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/domains/webhooks"
)

type ProcessorFactory interface {
	CreateMessageProcessor(connection connecting.Connection, tenantKey *tenants.TenantKey) (MessageProcessor, error)
}

type MessageProcessorFactory struct {
	webHooksService       webhooks.Service
	encryptionService     encryption.Service
	deliveryQueuesService delivery_queues.Service
	metrics               abstraction.MetricsObserver
}

func CreateMessageProcessorFactory(webHooksService webhooks.Service, encryptionService encryption.Service, deliveryQueuesService delivery_queues.Service, metricsObserver abstraction.MetricsObserver) MessageProcessorFactory {
	return MessageProcessorFactory{
		webHooksService:       webHooksService,
		encryptionService:     encryptionService,
		deliveryQueuesService: deliveryQueuesService,
		metrics:               metricsObserver,
	}
}

func (m MessageProcessorFactory) CreateMessageProcessor(connection connecting.Connection, tenantKey *tenants.TenantKey) (MessageProcessor, error) {
	switch connection.ProcessorType {
	case MessageProcessorTypeCDC:
		return &MessageProcessorCDC{DefaultMessageProcessor{
			tenantKey:             tenantKey,
			connection:            connection,
			webHooksService:       m.webHooksService,
			encryptionService:     m.encryptionService,
			deliveryQueuesService: m.deliveryQueuesService,
			metrics:               m.metrics,
		}}, nil
	case MessageProcessorTypeOdin:
		return &MessageProcessorOdin{DefaultMessageProcessor{
			tenantKey:             tenantKey,
			connection:            connection,
			webHooksService:       m.webHooksService,
			encryptionService:     m.encryptionService,
			deliveryQueuesService: m.deliveryQueuesService,
			metrics:               m.metrics,
		}}, nil
	case MessageProcessorTypeVNext:
		return &MessageProcessorVNext{DefaultMessageProcessor{
			tenantKey:             tenantKey,
			connection:            connection,
			webHooksService:       m.webHooksService,
			encryptionService:     m.encryptionService,
			deliveryQueuesService: m.deliveryQueuesService,
			metrics:               m.metrics,
		}}, nil
	default:
		return nil, errors.New("invalid processor type")
	}
}
