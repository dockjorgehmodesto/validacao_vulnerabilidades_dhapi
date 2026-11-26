package connectors

import (
	"context"
	"strconv"
	"sync"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/message_processors"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"conductor.com.br/pier/hooks/pkg/helpers/timer"
	billing "dev.azure.com/conductortech/billing/_git/middleware-golang.git"
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/Azure/go-amqp"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type serviceBusConnector struct {
	Connection        connecting.Connection
	tenantKey         *tenants.TenantKey     // Tenant used for the Connection
	ConnectedAt       time.Time              // When the application connected to the queue/topic
	subscription      ServiceBusSubscription // The active subscription for the connection. This is used for closing a connection if necessary
	context           context.Context        // Context is also used for handling/closing a connection
	contextCancelFunc context.CancelFunc     // Function for closing context
	mutex             sync.Mutex             // Makes sure that a connection is not changed when an action is being executed
	OnClose           func(key string)       // Function for disconnecting from manager

	forceClose      bool // Variable for forcing a connection to get out of its loop
	isDisconnecting bool // Won't allow Disconnect to be called twice
	processor       message_processors.MessageProcessor
}

func (s *serviceBusConnector) GetProcessorType() string {
	return s.Connection.ProcessorType
}

func (s *serviceBusConnector) Listen() {
	// TODO unnecessary goroutine
	go s.handleListening()
}

func (s *serviceBusConnector) setupSubscription() (ServiceBusSubscription, error) {
	namespace, e := servicebus.NewNamespace(s.getNamespaceWithConnectionString())
	if e != nil {
		return nil, errors.Wrap(e, "could not connect to namespace")
	}

	topic, e := namespace.NewTopic(s.Connection.Topic)
	if e != nil {
		return nil, errors.Wrap(e, "could not start topic")
	}

	subscription, e := topic.NewSubscription(s.Connection.SubscriptionName, servicebus.SubscriptionWithPrefetchCount(viper.GetUint32("amqp.servicebusPrefetch")))
	if e != nil {
		return nil, errors.Wrap(e, "could not create subscription")
	}
	return subscription, nil
}

// handleListening uses the subscription client to call the Receive loop function
// from azure SDK and passes a callback function responsible for processing the message.
// It should not stop and needs to block this process execution. That's why it is called
// using a go routine.
func (s *serviceBusConnector) handleListening() {

	defer s.Disconnect()

	// Message concurrency logic
	var concurrentHandler servicebus.HandlerFunc
	var msgConcurrencyChan chan *servicebus.Message
	defer func() {
		if msgConcurrencyChan != nil {
			close(msgConcurrencyChan)
		}
	}()
	{
		var concurrentNum = global.Viper.GetInt("amqp.workerCount.servicebusConnector")
		msgConcurrencyChan = make(chan *servicebus.Message, concurrentNum)
		concurrentHandler = func(ctx context.Context, msg *servicebus.Message) error {
			msgConcurrencyChan <- msg
			return nil
		}

		for i := 0; i < concurrentNum; i++ {
			go func() {
				for msg := range msgConcurrencyChan {
					ctx, cancelContext := context.WithCancel(context.Background())
					_ = s.processMessage(ctx, msg, message_processors.MessageSourceDefault)
					cancelContext()
				}
			}()
		}
	}

	// Logic for keeping a connection alive
	{
		for true {
			subscription, err := s.setupSubscription()
			if err == nil {

				s.context, s.contextCancelFunc = context.WithCancel(context.Background())
				err = subscription.Receive(
					s.context,
					global.BillingInterceptors.ServicebusInterceptor(
						s.GetIssuerId(),
						s.GetEventType(),
						billing.ServiceBusConnectionInfo{
							Topic:        s.Connection.Topic,
							Subscription: s.Connection.SubscriptionName,
						},
						concurrentHandler,
					),
				)

				if err == nil {
					err = errors.New("closed without errors")
				}

			}
			splunk.NewAMQPLog().SendCustomError("subscription reset", err, splunk.MessagingErrorMetadata{
				IssuerId:     s.Connection.IssuerId,
				EventType:    s.Connection.EventType,
				Subscription: s.Connection.SubscriptionName,
				ConnectionId: s.Connection.ID,
				AmqpSource:   connecting.AMQPBrokerServiceBus,
			})
			if s.forceClose {
				break
			}
			time.Sleep(viper.GetDuration("amqp.reconnectTimeoutSeconds") * time.Second)

		}

	}
	splunk.NewAMQPLog().SendCustomError("subscription closed", errors.New("forced to close"), splunk.MessagingErrorMetadata{
		IssuerId:     s.Connection.IssuerId,
		EventType:    s.Connection.EventType,
		Subscription: s.Connection.SubscriptionName,
		ConnectionId: s.Connection.ID,
		AmqpSource:   connecting.AMQPBrokerServiceBus,
	})
}

func getServiceBusDeliveryGeneratedAtHeader(msg *servicebus.Message) *time.Time {
	if generatedAtInterface, ok := msg.GetKeyValues()[message_processors.HeaderGeneratedAt]; ok {
		return safelyConvertStringInterfaceToTime(generatedAtInterface)
	}
	return nil
}

func getServiceBusDeliveryPublishAtHeader(msg *servicebus.Message) *time.Time {
	if publishAtInterface, ok := msg.GetKeyValues()[message_processors.HeaderPublishAt]; ok {
		return safelyConvertStringInterfaceToTime(publishAtInterface)
	}
	return nil
}

func getServiceBusDeliveryPrimaryKeyHeader(msg *servicebus.Message) uint64 {
	if primaryKeyInterface, ok := msg.GetKeyValues()[message_processors.HeaderPrimaryKey]; ok {
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

func getServiceBusDeliveryOperationHeader(msg *servicebus.Message) string {
	if operationInterface, ok := msg.GetKeyValues()[message_processors.HeaderOperation]; ok {
		switch operationInterface.(type) {
		case string:
			return operationInterface.(string)
		default:
			return ""
		}
	}
	return ""
}

func safelyConvertStringInterfaceToTime(stringAsInterface interface{}) *time.Time {
	if typedString, ok := stringAsInterface.(string); ok {
		parse, err := time.Parse(time.RFC3339, typedString)
		if err == nil {
			return &parse
		}
	}
	return nil
}

// Creates contexts, handles logs and messages destinations
func (s *serviceBusConnector) processMessage(ctx context.Context, msg *servicebus.Message, msgSource string) error {
	var e error
	tracerSpan, ctx := tracer.StartSpanFromContext(
		ctx,
		"servicebusread",
		tracer.Tag("issuerId", s.Connection.IssuerId),
		tracer.Tag("eventType", s.GetEventType()),
		tracer.Tag("requestId", msg.ID),
	)
	defer tracerSpan.Finish()

	groupId := ""
	if mapInterfaceResult, ok := msg.GetKeyValues()[message_processors.HeaderGroupId]; ok {
		if strValue, ok := mapInterfaceResult.(string); ok {
			groupId = strValue
		}
	}
	ctx = timer.AppendTrackerToContext(ctx)

	_, e = s.processor.Process(
		ctx,
		message_processors.Message{
			Id:            msg.ID,
			GroupId:       groupId,
			Data:          string(msg.Data),
			Source:        msgSource,
			DeliveryCount: int(msg.DeliveryCount),
			GeneratedAt:   getServiceBusDeliveryGeneratedAtHeader(msg),
			PrimaryKey:    getServiceBusDeliveryPrimaryKeyHeader(msg),
			Operation:     getServiceBusDeliveryOperationHeader(msg),
			PublishAt:     getServiceBusDeliveryPublishAtHeader(msg),
		},
	)

	if shouldRetryMessage(e) {
		return msg.Abandon(ctx)
	}
	return msg.Complete(ctx)
}

func (s *serviceBusConnector) Disconnect() error {
	if s.isDisconnecting {
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.isDisconnecting = true
	defer func() { s.isDisconnecting = false }()

	s.forceClose = true
	if s.subscription != nil {
		e := s.subscription.Close(s.context)
		if e != nil {
			splunk.NewAMQPLog().SendCustomError(
				"error executing subscription.Close() ",
				e,
				splunk.MessagingErrorMetadata{
					IssuerId:     s.Connection.IssuerId,
					EventType:    s.Connection.EventType,
					Subscription: s.Connection.SubscriptionName,
					ConnectionId: s.Connection.ID,
					AmqpSource:   connecting.AMQPBrokerServiceBus,
				},
			)
			if _, isDetachError := e.(*amqp.DetachError); !isDetachError {
				return e
			}
		}
	}

	if s.contextCancelFunc != nil {
		s.contextCancelFunc()
	}

	if s.OnClose != nil {
		s.OnClose(s.Connection.ID)
	}
	return nil
}

func (s *serviceBusConnector) Test() error {
	namespace, e := servicebus.NewNamespace(s.getNamespaceWithConnectionString())
	if e != nil {
		return errors.New("could not connect to namespace: " + e.Error())
	}

	topic, e := namespace.NewTopic(s.Connection.Topic)
	if e != nil {
		return errors.New("could not start topic: " + e.Error())
	}

	s.subscription, e = topic.NewSubscription(s.Connection.SubscriptionName)
	if e != nil {
		return errors.New("could not start subscription: " + e.Error())
	}

	context.Background()
	_, e = s.subscription.PeekOne(context.Background())
	if e != nil {
		if _, ok := e.(servicebus.ErrNoMessages); !ok {
			return errors.Wrap(e, "could not peek a message")
		}
	}
	return nil
}

func (s *serviceBusConnector) GetEventType() string {
	return s.Connection.EventType
}

func (s *serviceBusConnector) GetId() string {
	return s.Connection.ID
}

func (s *serviceBusConnector) GetIssuerId() int {
	return s.Connection.IssuerId
}

func (s *serviceBusConnector) SetOnClose(onclose func(key string)) {
	s.OnClose = onclose
}

func (s *serviceBusConnector) GetConnection() *connecting.Connection {
	return &s.Connection
}

func (s *serviceBusConnector) getNamespaceWithConnectionString() servicebus.NamespaceOption {
	return servicebus.NamespaceWithConnectionString(s.Connection.ConnectionString)
}
