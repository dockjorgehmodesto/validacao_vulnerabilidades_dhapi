package domains

import (
	"conductor.com.br/pier/hooks/pkg/domains/abstraction"
	"conductor.com.br/pier/hooks/pkg/domains/amqp"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/delivery"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/message_processors"
	"conductor.com.br/pier/hooks/pkg/domains/delivery_queues"
	"conductor.com.br/pier/hooks/pkg/domains/encryption"
	"conductor.com.br/pier/hooks/pkg/domains/event_controller"
	"conductor.com.br/pier/hooks/pkg/domains/issuer_events"
	"conductor.com.br/pier/hooks/pkg/domains/manager"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/domains/webhooks"
	cache2 "conductor.com.br/pier/hooks/pkg/helpers/cache"

	"github.com/eko/gocache/cache"
	"github.com/spf13/viper"
)

type CacheContext struct {
	Manager   cache.CacheInterface
	Marshaler cache2.MarshalerInterface
}

type repositories struct {
	encryption      encryption.Repository
	serviceBus      amqp.Repository
	webHook         webhooks.Repository
	webHookHeader   webhooks.RepositoryHeader
	issuerEvents    issuer_events.Repository
	eventController event_controller.Repository
	deliveryQueues  delivery_queues.Repository
}

type ServiceContainer struct {
	Encryption       encryption.Service
	AMQP             amqp.Service
	WebHook          webhooks.Service
	Tenant           tenants.Service
	IssuerEvents     issuer_events.Service
	WebHookDeliverer []delivery.Deliverer
	Management       manager.Service
	EventController  event_controller.Controller
	DeliveryQueues   delivery_queues.Service
}

func startRepositories(cacheMarshaler cache2.MarshalerInterface) repositories {
	return repositories{
		encryption:      encryption.NewCacheRepository(encryption.NewRepository(), cacheMarshaler),
		serviceBus:      amqp.NewRepository(),
		webHook:         webhooks.NewCacheRepository(webhooks.NewRepository(), cacheMarshaler),
		issuerEvents:    issuer_events.NewRepository(),
		eventController: event_controller.NewRepository(),
		deliveryQueues:  delivery_queues.NewCacheRepository(delivery_queues.NewRepository(), cacheMarshaler),
		webHookHeader:   webhooks.NewRepositoryHeader(),
	}
}

func Container(cache CacheContext, metricsObserver abstraction.MetricsObserver) ServiceContainer {
	repositories := startRepositories(cache.Marshaler)

	switcher := event_controller.NewSwitcher(cache.Manager)
	eventController := event_controller.NewEventController(switcher, repositories.eventController)

	encryptionService := encryption.NewService(repositories.encryption, tenants.TenantService)
	issuerEventsService := issuer_events.NewService(repositories.issuerEvents)
	webHookService := webhooks.NewService(repositories.webHook, encryptionService, issuerEventsService, eventController, repositories.webHookHeader)
	deliveryQueuesService := delivery_queues.NewService(repositories.deliveryQueues)
	processorFactory := message_processors.CreateMessageProcessorFactory(webHookService, encryptionService, deliveryQueuesService, metricsObserver)

	delivererWebhook := buildDeliveriesWebhook(metricsObserver)

	return ServiceContainer{
		Encryption:       encryptionService,
		WebHook:          webHookService,
		AMQP:             amqp.NewService(repositories.serviceBus, amqp.NewConnectionManager(), processorFactory, eventController),
		Tenant:           tenants.TenantService,
		IssuerEvents:     issuerEventsService,
		WebHookDeliverer: delivererWebhook,
		Management:       manager.NewManagementService(),
		EventController:  eventController,
		DeliveryQueues:   deliveryQueuesService,
	}
}

func buildDeliveriesWebhook(metricsObserver abstraction.MetricsObserver) []delivery.Deliverer {
	var webhookDelivery []delivery.Deliverer

	delivererRetry := delivery.NewWebHookDeliverer(metricsObserver)
	delivererRetry.BuildDelivery(delivery.CustomConfigs{
		QueueName:          delivery.QueueNameRetry,
		QueueKey:           delivery.QueueKeyRetry,
		ExchangeName:       delivery.ExchangeNameRetry,
		IsRetry:            true,
		RabbitRetryChannel: delivererRetry.GetAmqpChannel(),
	})
	webhookDelivery = append(webhookDelivery, delivererRetry)

	if !viper.GetBool("features.disableWebhookDelivery") {
		deliverer := delivery.NewWebHookDeliverer(metricsObserver)
		deliverer.BuildDelivery(delivery.CustomConfigs{
			QueueName:          delivery.QueueName,
			QueueKey:           delivery.QueueKey,
			ExchangeName:       delivery.ExchangeName,
			IsRetry:            false,
			RabbitRetryChannel: delivererRetry.GetAmqpChannel(),
		})
		webhookDelivery = append(webhookDelivery, deliverer)
	}
	return webhookDelivery
}
