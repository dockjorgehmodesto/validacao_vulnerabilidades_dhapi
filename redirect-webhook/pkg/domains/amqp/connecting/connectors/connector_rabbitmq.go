package connectors

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/message_processors"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"conductor.com.br/pier/hooks/pkg/helpers/timer"
	billing "dev.azure.com/conductortech/billing/_git/middleware-golang.git"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type rabbitMQConnector struct {
	Connection  connecting.Connection
	tenantKey   *tenants.TenantKey // Tenant used for the Connection
	ConnectedAt time.Time          // When the application connected to the queue/topic
	OnClose     func(key string)   // Function for disconnecting from manager

	rabbitConnection *amqp.Connection
	rabbitChannel    *amqp.Channel

	reconnectMutex sync.Mutex
	isReconnecting bool

	processor       message_processors.MessageProcessor
	isDisconnecting bool
	disconnectMutex sync.Mutex
}

func (r *rabbitMQConnector) GetProcessorType() string {
	return r.Connection.ProcessorType
}

func (r *rabbitMQConnector) Test() error {
	var err error
	r.rabbitConnection, err = amqp.DialConfig(
		r.Connection.ConnectionString,
		amqp.Config{
			Heartbeat: viper.GetDuration("amqp.dynamicRabbit.defaultHeartbeatSeconds") * time.Second,
			Locale:    viper.GetString("amqp.dynamicRabbit.defaultLocale"),
			Vhost:     r.Connection.VirtualHost,
		},
	)

	if err != nil {
		return err
	}

	r.rabbitChannel, err = r.rabbitConnection.Channel()
	if err != nil {
		return err
	}

	err = r.rabbitChannel.ExchangeDeclare(
		r.Connection.ExchangeName,
		viper.GetString("amqp.dynamicRabbit.exchangesDeclaration.kind"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.durable"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.autoDelete"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.internal"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.noWait"),
		nil,
	)

	if err != nil {
		return err
	}

	_, err = r.rabbitChannel.QueueDeclare(
		r.Connection.QueueName,
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.durable"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.autoDelete"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.exclusive"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.noWait"),
		nil,
	)

	if err != nil {
		return err
	}

	err = r.rabbitChannel.QueueBind(
		r.Connection.QueueName,
		r.Connection.QueueKey,
		r.Connection.ExchangeName,
		viper.GetBool("amqp.dynamicRabbit.queueBindings.noWait"),
		nil,
	)

	if err != nil {
		return err
	}
	hostname, _ := os.Hostname()
	_, err = r.rabbitChannel.Consume(
		r.Connection.QueueName,
		hostname,
		viper.GetBool("amqp.dynamicRabbit.consumer.autoAck"),
		viper.GetBool("amqp.dynamicRabbit.consumer.exclusive"),
		viper.GetBool("amqp.dynamicRabbit.consumer.noLocal"),
		viper.GetBool("amqp.dynamicRabbit.consumer.noWait"),
		nil,
	)

	_ = r.rabbitChannel.Close()
	_ = r.rabbitConnection.Close()
	return err
}

func (r *rabbitMQConnector) Listen() {
	r.panicHandler(r.rabbitMQConnectLogic)
}

func (r *rabbitMQConnector) panicHandler(criticalFunction func()) {
	defer func() {
		if fatalError := recover(); fatalError != nil {
			if val, ok := fatalError.(error); ok {
				r.handleError("fatal error", val, true)
			} else {
				logging.Logger.Error(fmt.Sprintf("[CONNECTOR] %v", fatalError))
			}
		}
	}()
	criticalFunction()
}

func (r *rabbitMQConnector) logToBilling(issuerId int, message amqp.Delivery) {
	global.BillingLogger.RabbitMQMessage(
		issuerId,
		r.Connection.EventType,
		billing.RabbitMQConnectionInfo{
			QueueKey:     r.Connection.QueueKey,
			QueueName:    r.Connection.QueueName,
			ExchangeName: r.Connection.ExchangeName,
		},
		message,
	)
}

func (r *rabbitMQConnector) rabbitMQConnectLogic() {
	var err error
	//
	// 1 Create connection
	//
	r.rabbitConnection, err = amqp.DialConfig(
		r.Connection.ConnectionString,
		amqp.Config{
			Heartbeat: viper.GetDuration("amqp.dynamicRabbit.defaultHeartbeatSeconds") * time.Second,
			Locale:    viper.GetString("amqp.dynamicRabbit.defaultLocale"),
			Vhost:     r.Connection.VirtualHost,
		},
	)

	if err != nil {
		r.handleError("could not open connection", err, true)
		return
	}

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

	// Listens for closed channels
	go r.closeListener("channel closed", r.rabbitChannel.NotifyClose)

	err = r.rabbitChannel.Qos(viper.GetInt("rabbit.qos.prefetchCount"), viper.GetInt("rabbit.qos.prefetchSize"), true)
	if err != nil {
		r.handleError("could not apply qos", err, true)
		return
	}

	//
	// 3 Declare exchange
	//
	err = r.rabbitChannel.ExchangeDeclare(
		r.Connection.ExchangeName,
		viper.GetString("amqp.dynamicRabbit.exchangesDeclaration.kind"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.durable"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.autoDelete"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.internal"),
		viper.GetBool("amqp.dynamicRabbit.exchangesDeclaration.noWait"),
		nil,
	)

	if err != nil {
		r.handleError("could not declare exchange", err, true)
		return
	}

	//
	// 4 Declare queue
	//
	_, err = r.rabbitChannel.QueueDeclare(
		r.Connection.QueueName,
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.durable"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.autoDelete"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.exclusive"),
		viper.GetBool("amqp.dynamicRabbit.queuesDeclaration.noWait"),
		nil,
	)

	if err != nil {
		r.handleError("could not declare queue", err, true)
		return
	}

	//
	// 5 Bind queue to exchange
	//
	err = r.rabbitChannel.QueueBind(
		r.Connection.QueueName,
		r.Connection.QueueKey,
		r.Connection.ExchangeName,
		viper.GetBool("amqp.dynamicRabbit.queueBindings.noWait"),
		nil,
	)

	if err != nil {
		r.handleError("could not bind queue to exchange", err, true)
		return
	}

	//
	// 6 Init listeners
	//
	hostname, _ := os.Hostname()
	deliveries, err := r.rabbitChannel.Consume(
		r.Connection.QueueName,
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

	helpers.ReadRabbitMsgsInGoroutines(global.Viper.GetInt("amqp.workerCount.rabbitMQConnector"), r.receiveRabbitMessage, deliveries)
}

func (r *rabbitMQConnector) receiveRabbitMessage(ctx context.Context, cancelContext context.CancelFunc, message amqp.Delivery) {
	e := func(delivery amqp.Delivery) error {
		defer cancelContext()

		if delivery.MessageId == "" {
			delivery.MessageId = uuid.New().String()
		}
		ctx = timer.AppendTrackerToContext(ctx)

		issuerId, err := r.processor.Process(
			ctx,
			message_processors.Message{
				Id:            delivery.MessageId,
				Data:          string(delivery.Body),
				Source:        message_processors.MessageSourceDefault,
				DeliveryCount: 0, // TODO needs to be reviewed, as Rabbit does not natively implement retry counters
				GroupId:       getRabbitDeliveryGroupId(message),
				GeneratedAt:   getRabbitDeliveryGeneratedAtHeader(message),
				PrimaryKey:    getRabbitDeliveryPrimaryKeyHeader(message),
				Operation:     getRabbitDeliveryOperationHeader(message),
				PublishAt:     getRabbitDeliveryPublishAtHeader(message),
			},
		)

		r.logToBilling(issuerId, message)

		if shouldRetryMessage(err) {
			return delivery.Nack(false, true)
		}
		return delivery.Ack(false)
	}(message)
	if e != nil {
		r.handleError("could not read message", e, false)
	}
}

func getRabbitDeliveryGroupId(message amqp.Delivery) string {
	groupId := ""
	if mapInterfaceResult, ok := message.Headers[message_processors.HeaderGroupId]; ok {
		if strValue, ok := mapInterfaceResult.(string); ok {
			groupId = strValue
		}
	}

	return groupId
}

func getRabbitDeliverySentAtHeader(message amqp.Delivery) *time.Time {
	if sentAtInterface, ok := message.Headers[message_processors.HeaderSentAt]; ok {
		return safelyConvertStringInterfaceToTime(sentAtInterface)
	}
	return nil
}

func getRabbitDeliveryGeneratedAtHeader(msg amqp.Delivery) *time.Time {
	if generatedAtInterface, ok := msg.Headers[message_processors.HeaderGeneratedAt]; ok {
		return safelyConvertStringInterfaceToTime(generatedAtInterface)
	}
	if sentAtInterface, ok := msg.Headers[message_processors.HeaderSentAt]; ok {
		return safelyConvertStringInterfaceToTime(sentAtInterface)
	}
	return nil
}

func getRabbitDeliveryPublishAtHeader(msg amqp.Delivery) *time.Time {
	if publishAtInterface, ok := msg.Headers[message_processors.HeaderPublishAt]; ok {
		return safelyConvertStringInterfaceToTime(publishAtInterface)
	}
	return nil
}

func getRabbitDeliveryPrimaryKeyHeader(msg amqp.Delivery) uint64 {
	if primaryKeyInterface, ok := msg.Headers[message_processors.HeaderPrimaryKey]; ok {
		switch primaryKeyInterface.(type) {
		case int:
			return uint64(primaryKeyInterface.(int))
		case int64:
			return uint64(primaryKeyInterface.(int64))
		case uint64:
			return primaryKeyInterface.(uint64)
		case string:
			atoi, err := strconv.Atoi(primaryKeyInterface.(string))
			if err != nil {
				return 0
			}
			return uint64(atoi)
		default:
			return 0
		}
	}
	return 0
}

func getRabbitDeliveryOperationHeader(msg amqp.Delivery) string {
	if operationInterface, ok := msg.Headers[message_processors.HeaderOperation]; ok {
		switch operationInterface.(type) {
		case string:
			return operationInterface.(string)
		default:
			return ""
		}
	}
	return ""
}

func (r *rabbitMQConnector) closeListener(errMsg string, notifier func(chan *amqp.Error) chan *amqp.Error) {
	func() {
		errChan := make(chan *amqp.Error)
		for {
			reason, ok := <-notifier(errChan)
			if ok {
				r.handleError(errMsg, reason, true)
				return
			}
		}
	}()
}

func (r *rabbitMQConnector) handleError(msg string, err error, shouldRestart bool) {

	splunk.NewAMQPLog().SendCustomError(
		msg,
		err,
		splunk.MessagingErrorMetadata{
			IssuerId:     r.Connection.IssuerId,
			EventType:    r.Connection.EventType,
			ConnectionId: r.Connection.ID,
			AmqpSource:   connecting.AMQPBrokerRabbitMQ,
		},
	)

	if shouldRestart {
		r.restart()
	}
}

func (r *rabbitMQConnector) setReconnecting(b bool) {
	r.reconnectMutex.Lock()
	defer r.reconnectMutex.Unlock()
	r.isReconnecting = b
}

func (r *rabbitMQConnector) restart() {
	if !r.isReconnecting {
		r.setReconnecting(true)
		defer r.setReconnecting(false)
		r.closeConnectionAndChannel()

		timeout := viper.GetDuration("amqp.dynamicRabbit.restartTimeoutSeconds")
		time.Sleep(timeout * time.Second)

		splunkLog := splunk.AMQPLog{}
		splunkLog.InitRabbitMQLog("", r.GetEventType(), r.GetIssuerId(), 0, r.Connection.ID)
		splunkLog.SendCustomMessage(splunk.LevelWarn, "restarting")
		r.Listen()
	}
}

func (r *rabbitMQConnector) closeConnectionAndChannel() {
	if r.rabbitChannel != nil {
		_ = r.rabbitChannel.Close()
	}

	if r.rabbitConnection != nil {
		_ = r.rabbitConnection.Close()
	}
}

func (r *rabbitMQConnector) Disconnect() error {
	if r.isDisconnecting {
		return nil
	}

	r.disconnectMutex.Lock()
	defer r.disconnectMutex.Unlock()

	r.isDisconnecting = true
	defer func() { r.isDisconnecting = false }()

	r.closeConnectionAndChannel()
	r.OnClose(r.Connection.ID)
	return nil
}

func (r *rabbitMQConnector) GetEventType() string {
	return r.Connection.EventType
}

func (r *rabbitMQConnector) GetId() string {
	return r.Connection.ID
}

func (r *rabbitMQConnector) GetIssuerId() int {
	return r.Connection.IssuerId
}

func (r *rabbitMQConnector) SetOnClose(onclose func(key string)) {
	r.OnClose = onclose
}

func (r *rabbitMQConnector) GetConnection() *connecting.Connection {
	return &r.Connection
}
