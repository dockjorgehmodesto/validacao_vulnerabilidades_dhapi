package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/abstraction"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	rabbitmq_x_delay "conductor.com.br/pier/hooks/pkg/helpers/retry"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"

	billing "dev.azure.com/conductortech/billing/_git/middleware-golang.git"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var (
	histogramProcessDelay = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "message_processing_delay_ms",
		Help: "The delay in milliseconds between sending and receiving a message from RabbitMQ",
	}, []string{"issuerId", "eventType"})
	histogramOriginDelay = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "message_origin_delay_ms",
		Help: "The delay in milliseconds between the event being generated on its origin and sending it to the issuer",
	}, []string{"issuerId", "eventType"})
	ExchangeName      = viper.GetString("rabbit.defaultExchangeName")
	QueueName         = viper.GetString("rabbit.defaultQueueName")
	QueueKey          = viper.GetString("rabbit.defaultQueueKey")
	ExchangeNameRetry = viper.GetString("rabbit.defaultRetryExchangeName")
	QueueNameRetry    = viper.GetString("rabbit.defaultRetryQueueName")
	QueueKeyRetry     = viper.GetString("rabbit.defaultRetryQueueKey")
)

type Deliverer interface {
	StartListening()
	BuildDelivery(deliveryConfig CustomConfigs)
	GetAmqpChannel() **amqp.Channel
}

type WebHookDeliverer struct {
	rabbitConnection *amqp.Connection
	rabbitChannel    *amqp.Channel
	reconnectMutex   sync.Mutex
	isReconnecting   bool
	httpClient       *http.Client

	rabbitRetryChannel **amqp.Channel

	queueName    string
	queueKey     string
	exchangeName string

	isRetry bool
	metrics abstraction.MetricsObserver
}

type CustomConfigs struct {
	QueueName    string
	QueueKey     string
	ExchangeName string

	RabbitRetryChannel **amqp.Channel
	IsRetry            bool
}

func NewWebHookDeliverer(metrics abstraction.MetricsObserver) Deliverer {
	return &WebHookDeliverer{
		reconnectMutex: sync.Mutex{},
		isReconnecting: false,
		httpClient: &http.Client{
			// Values from http.DefaultTransport
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   viper.GetDuration("http.transport.dialContext.timeoutSecs") * time.Second,
					KeepAlive: viper.GetDuration("http.transport.dialContext.keepAliveSecs") * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          viper.GetInt("http.transport.maxIdleConns"),
				IdleConnTimeout:       viper.GetDuration("http.transport.idleConnTimeoutSecs") * time.Second,
				TLSHandshakeTimeout:   viper.GetDuration("http.transport.tlsHandshakeTimeoutSecs") * time.Second,
				ExpectContinueTimeout: viper.GetDuration("http.transport.expectContinueTimeoutSecs") * time.Second,
			},
			Timeout: viper.GetDuration("http.timeoutSecs") * time.Second,
		},
		metrics: metrics,
	}
}

func (r *WebHookDeliverer) BuildDelivery(deliveryConfig CustomConfigs) {
	r.queueName = deliveryConfig.QueueName
	r.queueKey = deliveryConfig.QueueKey
	r.exchangeName = deliveryConfig.ExchangeName
	r.rabbitRetryChannel = deliveryConfig.RabbitRetryChannel
	r.isRetry = deliveryConfig.IsRetry
}

func (r *WebHookDeliverer) GetAmqpChannel() **amqp.Channel {
	return &r.rabbitChannel
}

func (r *WebHookDeliverer) logToBilling(ctx context.Context, webHookData WebHookMessageDelivery, statusCode int, bytesThroughput int64) {
	tracerSpan, ctx := tracer.StartSpanFromContext(ctx, "logtobilling")
	defer tracerSpan.Finish()
	global.BillingLogger.WebHookDelivery(billing.WebHookMessage{
		StatusCode: statusCode,
		URL:        webHookData.URL,
		Headers:    webHookData.Headers,
		Metadata: billing.Metadata{
			MessageId:    webHookData.Metadata.MessageId,
			InternalId:   webHookData.Metadata.InternalId,
			KeyVaultId:   webHookData.Metadata.KeyVaultId,
			WebHookId:    webHookData.Metadata.WebHookId,
			ConnectionId: webHookData.Metadata.ConnectionId,
			EventType:    webHookData.Metadata.EventType,
			BrokerType:   webHookData.Metadata.BrokerType,
			Host:         webHookData.Metadata.Host,
			ProcessedAt:  webHookData.Metadata.ProcessedAt,
		},
		IssuerId:        webHookData.IssuerId,
		BytesThroughput: bytesThroughput,
	})
}

func (r *WebHookDeliverer) readMessageAndSendWebHook(ctx context.Context, cancelContext context.CancelFunc, amqpDelivery amqp.Delivery) {
	e := func(delivery amqp.Delivery) error {
		tracerSpan, ctx := tracer.StartSpanFromContext(ctx, "webhookdelivery")
		defer tracerSpan.Finish()
		defer cancelContext()

		webhook := WebHookMessageDelivery{}
		err := json.Unmarshal(delivery.Body, &webhook)
		if err != nil {
			splunk.NewDeliveryLog(0).WithError(err).WithMessage("could not unmarshal delivery").WithLevel(splunk.LevelError).Send(ctx, "could not unmarshal delivery")
			return delivery.Reject(false)
		}
		tracerSpan.SetTag("issuerId", webhook.IssuerId)
		tracerSpan.SetTag("eventType", webhook.Metadata.EventType)
		tracerSpan.SetTag("messageId", webhook.Metadata.MessageId)

		deliveryLog := r.createDeliveryLog(delivery, webhook)
		histogramProcessDelay.
			WithLabelValues(fmt.Sprintf("%d", webhook.IssuerId), webhook.Metadata.EventType).
			Observe(float64(time.Since(webhook.Metadata.ProcessedAt).Milliseconds()))
		defer histogramOriginDelay.
			WithLabelValues(fmt.Sprintf("%d", webhook.IssuerId), webhook.Metadata.EventType).
			Observe(float64(webhook.Metadata.GetOriginDelayMs()))

		if webhook.Headers == nil {
			webhook.Headers = make(map[string]string)
		}
		webhook.Headers["DataHubMessageDeliveryCount"] = strconv.Itoa(getDeliveryCount(delivery))
		webhook.Headers["EventType"] = webhook.Metadata.EventType
		response, err := helpers.HttpRequest(ctx, r.httpClient, webhook.URL, webhook.Headers, []byte(webhook.Payload), http.MethodPost)

		tags := map[string]interface{}{
			"event_name":     deliveryLog.Trace.Metadata.EventType,
			"issuer_id":      deliveryLog.Trace.IssuerId,
			"env":            deliveryLog.AppProfile,
			"service":        deliveryLog.Type,
			"delivery_count": deliveryLog.Trace.Metadata.DeliveryCount,
		}

		// TODO: investigar o motivo do campo OrgifinDelayMillis vir vazio
		if deliveryLog.Trace.Metadata.OriginDelayMillis != nil {
			r.metrics.RecordDuration("origin_delay_millis", float64(*deliveryLog.Trace.Metadata.OriginDelayMillis), tags)
		}

		if deliveryLog.Trace.Metadata.PublishDelayMillis != nil {
			r.metrics.RecordDuration("publish_delay_millis", float64(*deliveryLog.Trace.Metadata.PublishDelayMillis), tags)
		}

		if err != nil {
			deliveryLog.WithError(err).WithMessage("http client error").WithLevel(splunk.LevelError).Send(ctx, "http client error")
			r.reschedule(delivery, ctx)
			return delivery.Reject(false)
		}

		r.logToBilling(ctx, webhook, response.Status, response.TotalBytesThroughput)
		if response.Status > 499 {
			deliveryLog.WithResponse(response).WithMessage("invalid status code").WithLevel(splunk.LevelWarn).Send(ctx, "invalid status code")
			r.reschedule(delivery, ctx)
			return delivery.Reject(false)
		}

		deliveryLog.WithResponse(response).WithMessage("success").WithLevel(splunk.LevelInfo).Send(ctx, "event delivered to the client")
		return delivery.Ack(false)
	}(amqpDelivery)
	if e != nil {
		r.handleError("could not read message", e, false)
	}
}

func (r *WebHookDeliverer) createDeliveryLog(delivery amqp.Delivery, message WebHookMessageDelivery) *splunk.DeliveryLog {
	deliveryLog := splunk.NewDeliveryLog(message.IssuerId)

	deliveryLog.WithTrace(splunk.DeliveryTrace{
		MessageExternalId: message.Metadata.MessageId,
		MessageInternalId: message.Metadata.InternalId,
		GroupId:           message.Metadata.GroupId,
		Metadata: &splunk.DeliveryMetadata{
			EventType:          message.Metadata.EventType,
			WebHookId:          message.Metadata.WebHookId,
			DeliveryCount:      getDeliveryCount(delivery),
			ConnectionId:       message.Metadata.ConnectionId,
			KeyVaultId:         message.Metadata.KeyVaultId,
			ProcessedAt:        message.Metadata.ProcessedAt,
			ProcessDelayMillis: time.Since(message.Metadata.ProcessedAt).Milliseconds(),
			Host:               message.Metadata.Host,
			GeneratedAt:        message.Metadata.GeneratedAt,
			PrimaryKey:         message.Metadata.PrimaryKey,
			Operation:          message.Metadata.Operation,
			PublishAt:          message.Metadata.PublishAt,
		},
		Processor: &splunk.DeliveryProcessorTrace{
			InsertedOnDate: message.Processor.InsertedOnDate,
			DurationSec:    message.Processor.DurationSec,
			Durations:      message.Processor.Durations,
		},
	})
	return deliveryLog
}

func (r *WebHookDeliverer) reschedule(delivery amqp.Delivery, ctx context.Context) error {
	tracerSpan, ctx := tracer.StartSpanFromContext(ctx, "reschedulemessage")
	defer tracerSpan.Finish()
	return rabbitmq_x_delay.RescheduleMessage(delivery, &rabbitmq_x_delay.RescheduleOptions{
		Channel:      *r.rabbitRetryChannel,
		ExchangeName: ExchangeNameRetry,
		RoutingKey:   QueueKeyRetry,
	})
}

func (r *WebHookDeliverer) StartListening() {
	if !viper.GetBool("features.disableWebhookDelivery") || !viper.GetBool("features.disableWebhookDeliveryRetry") {
		logging.Logger.Info("[WEBHOOKS-CON] Starting")

		go r.rabbitMQConnectLogic()
	}
}

func (r *WebHookDeliverer) rabbitMQConnectLogic() {

	logging.Logger.Info(fmt.Sprintf("[WEBHOOKS-CON] Connecting queueName: %s queueKey: %s exchangeName: %s", r.queueName, r.queueKey, r.exchangeName))
	var err error
	//
	// 1 Create connection
	//
	r.rabbitConnection, err = amqp.DialConfig(
		viper.GetString("rabbit.dial"),
		amqp.Config{
			Heartbeat: viper.GetDuration("amqp.dynamicRabbit.defaultHeartbeatSeconds") * time.Second,
			Locale:    viper.GetString("amqp.dynamicRabbit.defaultLocale"),
			Vhost:     viper.GetString("rabbit.virtualHost"),
		},
	)

	if err != nil {
		r.handleError("could not open connection", err, true)
		return
	}
	logging.Logger.Info("[WEBHOOKS-CON] Created connection")

	// Listens for closed connections
	go r.closeListener("connection closed", r.rabbitConnection.NotifyClose)

	//
	// 2 Open channel
	//
	r.rabbitChannel, err = r.rabbitConnection.Channel()
	if err != nil {
		r.handleError("could not open channel", err, true)
		return
	}

	r.rabbitRetryChannel = &r.rabbitChannel

	err = r.rabbitChannel.Qos(viper.GetInt("rabbit.qos.prefetchCount"), viper.GetInt("rabbit.qos.prefetchSize"), true)
	if err != nil {
		r.handleError("could not apply qos", err, true)
		return
	}

	logging.Logger.Info("[WEBHOOKS-CON] Created channel")

	// Listens for closed channels
	go r.closeListener("channel closed", r.rabbitChannel.NotifyClose)

	//
	// 3 Declare exchange
	//
	err = rabbitmq_x_delay.CreateDefaultDelayExchange(r.exchangeName, viper.GetString("amqp.dynamicRabbit.exchangesDeclaration.kind"), r.rabbitChannel)

	if err != nil {
		r.handleError("could not declare exchange", err, true)
		return
	}

	logging.Logger.Info("[WEBHOOKS-CON] Created exchange")

	//
	// 4 Declare queue
	//
	_, err = r.rabbitChannel.QueueDeclare(
		r.queueName,
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.durable"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.autoDelete"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.exclusive"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.noWait"),
		amqp.Table{
			"x-message-ttl": viper.GetInt64("amqp.dynamicRabbit.queuesDeclaration.messageTimeToLiveMillis"),
		},
	)

	if err != nil {
		r.handleError("could not declare queue", err, true)
		return
	}

	logging.Logger.Info("[WEBHOOKS-CON] Created queue")

	//
	// 5 Bind queue to exchange
	//
	err = r.rabbitChannel.QueueBind(
		r.queueName,
		r.queueKey,
		r.exchangeName,
		viper.GetBool("amqp.dynamicRabbit.queueBindings.noWait"),
		nil,
	)

	if err != nil {
		r.handleError("could not bind queue to exchange", err, true)
		return
	}

	logging.Logger.Info("[WEBHOOKS-CON] Bound queue")

	isRetryApp := r.isRetry && !viper.GetBool("features.disableWebhookDeliveryRetry")
	if isRetryApp {
		r.startListeners()
	}

	if !r.isRetry && !viper.GetBool("features.disableWebhookDelivery") {
		r.startListeners()
	}
}

func (r *WebHookDeliverer) startListeners() {
	hostname, _ := os.Hostname()
	deliveries, err := r.rabbitChannel.Consume(
		r.queueName,
		hostname,
		viper.GetBool("amqp.dynamicRabbit.consumer.autoAck"),
		viper.GetBool("amqp.dynamicRabbit.consumer.exclusive"),
		viper.GetBool("amqp.dynamicRabbit.consumer.noLocal"),
		viper.GetBool("amqp.dynamicRabbit.consumer.noWait"),
		nil,
	)
	if err != nil {
		r.handleError("could not start consumer", err, true)
		return
	}
	logging.Logger.Info("[WEBHOOKS-CON] Started consuming", zap.String("queueName", r.queueName), zap.String("queueKey", r.queueKey), zap.String("exchangeName", r.exchangeName))

	helpers.ReadRabbitMsgsInGoroutines(global.Viper.GetInt("amqp.workerCount.webhookDeliverer"), r.readMessageAndSendWebHook, deliveries)
}

func (r *WebHookDeliverer) closeListener(errMsg string, notifier func(chan *amqp.Error) chan *amqp.Error) {
	errChan := make(chan *amqp.Error)
	for {
		if notifier != nil {
			reason, ok := <-notifier(errChan)
			if ok && reason != nil {
				r.handleError(errMsg, reason, true)
				return
			}
		}
	}
}

func (r *WebHookDeliverer) handleError(msg string, err error, shouldRestart bool) {
	logging.Logger.Error("[WEBHOOKS-CON] "+msg, zap.Error(err))
	r.logError(msg, err)

	if shouldRestart {
		r.restart()
	}
}

func (r *WebHookDeliverer) logError(msg string, err error) {
	errLog := splunk.NewAMQPLog()
	errLog.BrokerType = "WEBHOOKRABBITMQ"
	errLog.SendCustomError(
		msg,
		err,
		splunk.MessagingErrorMetadata{
			AmqpSource: connecting.AMQPBrokerRabbitMQ,
		},
	)
}

func (r *WebHookDeliverer) logWarn(msg string) {
	errLog := splunk.NewAMQPLog()
	errLog.BrokerType = "WEBHOOKRABBITMQ"
	errLog.SendCustomMessage(splunk.LevelWarn, msg)
}

func (r *WebHookDeliverer) setReconnecting(b bool) {
	r.reconnectMutex.Lock()
	defer r.reconnectMutex.Unlock()
	r.isReconnecting = b
}

func (r *WebHookDeliverer) restart() {
	if !r.isReconnecting {
		r.setReconnecting(true)
		defer r.setReconnecting(false)
		r.closeConnectionAndChannel()

		timeout := viper.GetDuration("amqp.dynamicRabbit.restartTimeoutSeconds")
		time.Sleep(timeout * time.Second)

		r.logWarn("restarting")
		r.StartListening()
	}
}

func (r *WebHookDeliverer) closeConnectionAndChannel() {
	if r.rabbitChannel != nil {
		_ = r.rabbitChannel.Close()
	}

	if r.rabbitConnection != nil {
		_ = r.rabbitConnection.Close()
	}
}
