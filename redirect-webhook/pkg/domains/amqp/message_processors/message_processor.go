package message_processors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"conductor.com.br/pier/hooks/pkg/domains/abstraction"
	"conductor.com.br/pier/hooks/pkg/domains/delivery_queues"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	connection2 "conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/delivery"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/rabbitmq"
	"conductor.com.br/pier/hooks/pkg/domains/encryption"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/domains/webhooks"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"conductor.com.br/pier/hooks/pkg/helpers/timer"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	MessageProcessorTypeCDC   = "CDC_MESSAGE_PROCESSOR"
	MessageProcessorTypeOdin  = "ODIN_MESSAGE_PROCESSOR"
	MessageProcessorTypeVNext = "VNEXT_MESSAGE_PROCESSOR"
)

type MessageProcessor interface {
	EncryptMessage(ctx context.Context, keyVault *encryption.KeyVault, message Message) (string, error)
	Send(ctx context.Context, delivery delivery.WebHookMessageDelivery, exchangeName string, queueKey string) error
	ListWebHooks(ctx context.Context, issuerId int, eventType string) ([]webhooks.WebHook, error)
	GetKeyVault(ctx context.Context, issuerId int) (*encryption.KeyVault, error)
	Process(ctx context.Context, message Message) (issuerId int, err error)
}

type DefaultMessageProcessor struct {
	tenantKey             *tenants.TenantKey
	connection            connection2.Connection
	webHooksService       webhooks.Service
	encryptionService     encryption.Service
	deliveryQueuesService delivery_queues.Service
	metrics               abstraction.MetricsObserver
}

func (d DefaultMessageProcessor) EncryptMessage(ctx context.Context, keyVault *encryption.KeyVault, message Message) (string, error) {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "encryptMessage")
	defer ddTrace.Finish()
	if keyVault == nil {
		return "", errors.New("nil keyvault on message encryption")
	}
	encryptedMessage, err := keyVault.EncryptWithAes(message.Data)
	if err != nil {
		return "", errors.Wrap(err, "could not encrypt webhook's payload")
	}
	return encryptedMessage, nil
}

func (d DefaultMessageProcessor) Send(ctx context.Context, webhookDelivery delivery.WebHookMessageDelivery, exchangeName string, queueKey string) error {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "sendProcessedMessage")
	defer ddTrace.Finish()
	defer timer.Track(ctx, splunk.DurationsKeyRabbitMQ, time.Now())
	msg, _ := json.Marshal(webhookDelivery)
	return rabbitmq.PublishSync(msg, exchangeName, queueKey, 1, 1, ctx)
}

func (d DefaultMessageProcessor) ListWebHooks(ctx context.Context, issuerId int, eventType string) ([]webhooks.WebHook, error) {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "fetchWebhooks")
	defer ddTrace.Finish()
	defer timer.Track(ctx, splunk.DurationsKeyWebhook, time.Now())
	relatedWebHooks, err := d.webHooksService.FindAllByTypeAndIssuerId(ctx, d.tenantKey, eventType, issuerId)

	if err != nil {
		return nil, err
	}
	if len(relatedWebHooks) == 0 {
		return nil, WebHookNotSetup
	}
	return relatedWebHooks, nil
}

func (d DefaultMessageProcessor) GetKeyVault(ctx context.Context, issuerId int) (*encryption.KeyVault, error) {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "fetchKeyVault")
	defer ddTrace.Finish()
	defer timer.Track(ctx, splunk.DurationsKeyKeyvault, time.Now())
	keyVault, err := d.encryptionService.FindMostRecentAndActiveByIssuerId(ctx, d.tenantKey, issuerId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, EncryptionNotFound
		} else {
			return nil, errors.Wrap(err, "error fetching keyvault")
		}
	}

	return &keyVault, err
}

func (d DefaultMessageProcessor) Process(context.Context, Message) (int, error) {
	panic("should not use the default implementation")
}

func (d DefaultMessageProcessor) process(ctx context.Context, message Message, issuerId int) (int, error) {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "messageprocessing")
	defer ddTrace.Finish()
	processorsLog := splunk.NewProcessorsLog(issuerId)

	// Treating scenarios where a single broker can dispatch messages for multiple events
	eventType := d.connection.EventType
	if message.EventType != "" {
		eventType = message.EventType
	}

	processorsLog.WithTrace(splunk.ProcessorTrace{
		MessageExternalId: message.Id,
		GroupId:           message.GroupId,
		EventType:         eventType,
		DeliveryCount:     message.DeliveryCount,
		ConnectionId:      d.connection.ID,
		BrokerType:        d.connection.Broker,
		GeneratedAt:       message.GeneratedAt,
		PrimaryKey:        message.PrimaryKey,
		Operation:         message.Operation,
		PublishAt:         message.PublishAt,
	})

	tags := map[string]interface{}{
		"event_name":     processorsLog.Trace.EventType,
		"issuer_id":      processorsLog.Trace.IssuerId,
		"env":            processorsLog.AppProfile,
		"service":        processorsLog.Type,
		"delivery_count": processorsLog.Trace.DeliveryCount,
	}

	// TODO: investigar o motivo do campo OrgifinDelayMillis vir vazio
	if processorsLog.Trace.OriginDelayMillis != nil {
		d.metrics.RecordDuration("origin_delay_millis", float64(*processorsLog.Trace.OriginDelayMillis), tags)
	}

	if processorsLog.Trace.PublishDelayMillis != nil {
		d.metrics.RecordDuration("publish_delay_millis", float64(*processorsLog.Trace.PublishDelayMillis), tags)
	}

	keyVault, err := d.GetKeyVault(ctx, issuerId)
	if err != nil {
		if err == EncryptionNotFound {
			processorsLog.WithError(err, splunk.LevelWarn)
		} else {
			processorsLog.WithError(err, splunk.LevelError)
		}
		processorsLog.Send(ctx, "error on get keyvault")
		return issuerId, err
	}
	processorsLog.Trace.EncryptionKeyId = keyVault.ID

	encryptedMessage, err := d.EncryptMessage(ctx, keyVault, message)
	if err != nil {
		processorsLog.WithError(err, splunk.LevelError).Send(ctx, "error encrypting event")
		return issuerId, err
	}

	webHooks, err := d.ListWebHooks(ctx, issuerId, eventType)
	if err != nil {
		if err == WebHookNotSetup {
			processorsLog.WithError(err, splunk.LevelWarn)
		} else {
			processorsLog.WithError(err, splunk.LevelError)
		}
		processorsLog.Send(ctx, "error searching for webhook")
		return issuerId, err
	}

	for _, webHook := range webHooks {
		logForCurrentMessage := *processorsLog
		logForCurrentMessage.Trace.MessageInternalId = uuid.New().String()
		logForCurrentMessage.Trace.WebHookId = webHook.ID

		parsedUrl, err := url.Parse(webHook.URL)

		if err != nil {
			logForCurrentMessage.WithError(err, splunk.LevelError)
			logging.Logger.Info(fmt.Sprintf("PARSED URL: %s", parsedUrl))
			return issuerId, err
		}

		if parsedUrl.Host == "" {
			logForCurrentMessage.WithError(err, splunk.LevelError)
			logging.Logger.Info("Parsed Host null")
			return issuerId, err
		}

		exchangeName := delivery.ExchangeName
		queueKey := delivery.QueueKey

		deliveryQueue, err := d.deliveryQueuesService.FindDeliveryQueue(d.tenantKey, eventType, issuerId)

		if err == nil {
			exchangeName = deliveryQueue.ExchangeName
			queueKey = deliveryQueue.QueueKey
		}

		if exchangeName == "" || queueKey == "" {
			exchangeName = global.Viper.GetString("rabbit.defaultExchangeName")
			queueKey = global.Viper.GetString("rabbit.defaultQueueKey")

			logForCurrentMessage.WithError(err, splunk.LevelError)

			logMessage := logForCurrentMessage
			logMessage.WithError(fmt.Errorf("the deliverySettings information is incorrect"), splunk.LevelError)
			logMessage.Send(ctx, "empty delivery settings")
		}

		err = d.Send(ctx, delivery.WebHookMessageDelivery{
			Payload: encryptedMessage,
			URL:     webHook.URL,
			Headers: webHook.WithHeaders(message.MakeHeaders()).GetHeaderMap(),
			Metadata: delivery.Metadata{
				MessageId:    message.Id,
				InternalId:   logForCurrentMessage.Trace.MessageInternalId,
				GroupId:      message.GroupId,
				KeyVaultId:   keyVault.ID,
				WebHookId:    webHook.ID,
				ConnectionId: d.connection.ID,
				EventType:    eventType,
				BrokerType:   d.connection.Broker,
				ProcessedAt:  time.Now(),
				GeneratedAt:  message.GeneratedAt,
				Host:         webHook.URL,
				PrimaryKey:   message.PrimaryKey,
				Operation:    message.Operation,
				PublishAt:    message.PublishAt,
			},
			Processor: delivery.ProcessorTrace{
				InsertedOnDate: logForCurrentMessage.Trace.InsertedOnDate,                            // logForCurrentMessage.Trace.SplunkTrace.InsertedOnDate,
				DurationSec:    time.Since(logForCurrentMessage.Trace.InsertedOnDate).Milliseconds(), //logForCurrentMessage.Trace.DurationSec,
				Durations:      timer.GetTracker(ctx),                                                // logForCurrentMessage.Trace.SplunkTrace.Durations,
			},
			IssuerId: issuerId,
		}, exchangeName, queueKey)

		if err != nil {
			logForCurrentMessage.WithError(err, splunk.LevelError)
		} else {
			logForCurrentMessage.Shipped = true
		}
		logForCurrentMessage.Send(ctx, "event sent to datahub delivery")

	}
	return issuerId, nil
}
