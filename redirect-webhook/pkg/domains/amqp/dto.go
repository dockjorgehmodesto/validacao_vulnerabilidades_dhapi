package amqp

import (
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/message_processors"
	"conductor.com.br/pier/hooks/pkg/helpers"
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/astaxie/beego/validation"
	"gopkg.in/jeevatkm/go-model.v1"
	"time"
)

type RabbitMQConnectionData struct {
	VirtualHost  string `json:"virtualHost"`
	ExchangeName string `json:"exchangeName"`
	QueueName    string `json:"queueName"`
	QueueKey     string `json:"queueKey"`
}

type ServiceBusConnectionData struct {
	Topic            string `json:"topic" valid:"Required"`
	SubscriptionName string `json:"subscriptionName" valid:"Required" example:"webhook"`
}

type AMQPPersist struct {
	AutoConnect   bool   `json:"autoConnect"`
	EventType     string `json:"eventType"`
	Broker        string `json:"broker" valid:"Required" example:"SERVICEBUS|RABBITMQ"`
	ProcessorType string `json:"processorType" valid:"Required" example:"CDC_MESSAGE_PROCESSOR|ODIN_MESSAGE_PROCESSOR|VNEXT_MESSAGE_PROCESSOR"`

	ConnectionString         string                    `json:"connectionString" valid:"Required"`
	RabbitMQConnectionData   *RabbitMQConnectionData   `json:"rabbitMqConnectionData"`
	ServiceBusConnectionData *ServiceBusConnectionData `json:"serviceBusConnectionData"`
}

func (s AMQPPersist) Valid(validation *validation.Validation) {
	if s.Broker != connecting.AMQPBrokerRabbitMQ && s.Broker != connecting.AMQPBrokerServiceBus {
		_ = validation.SetError("broker", "should be one of [RABBIMQ,SERVICEBUS]")
	}

	if s.ProcessorType != message_processors.MessageProcessorTypeCDC && s.ProcessorType != message_processors.MessageProcessorTypeOdin && s.ProcessorType != message_processors.MessageProcessorTypeVNext {
		_ = validation.SetError("processorType", "should be one of [CDC_MESSAGE_PROCESSOR,ODIN_MESSAGE_PROCESSOR,VNEXT_MESSAGE_PROCESSOR]")
	}

	if s.Broker == connecting.AMQPBrokerRabbitMQ && s.RabbitMQConnectionData == nil {
		_ = validation.SetError("rabbitMqConnectionData", "is required")
	}

	if s.Broker == connecting.AMQPBrokerServiceBus && s.ServiceBusConnectionData == nil {
		_ = validation.SetError("serviceBusConnectionData", "is required")
	}
}

func (s AMQPPersist) ConvertToEntity() connecting.Connection {
	if s.Broker == connecting.AMQPBrokerServiceBus {
		return connecting.Connection{
			ConnectionString: s.ConnectionString,
			Topic:            s.ServiceBusConnectionData.Topic,
			SubscriptionName: s.ServiceBusConnectionData.SubscriptionName,
			EventType:        s.EventType,
			Broker:           s.Broker,
			ProcessorType:    s.ProcessorType,
		}
	} else {
		return connecting.Connection{
			ConnectionString: s.ConnectionString,
			VirtualHost:      s.RabbitMQConnectionData.VirtualHost,
			ExchangeName:     s.RabbitMQConnectionData.ExchangeName,
			QueueName:        s.RabbitMQConnectionData.QueueName,
			QueueKey:         s.RabbitMQConnectionData.QueueKey,
			EventType:        s.EventType,
			Broker:           s.Broker,
			ProcessorType:    s.ProcessorType,
		}
	}
}

type AMQPUpdate struct {
	ConnectionString string `json:"connectionString" valid:"Required"`

	Broker                   string                    `json:"broker" valid:"Required" example:"SERVICEBUS|RABBITMQ"`
	ProcessorType            string                    `json:"processorType" valid:"Required" example:"CDC_MESSAGE_PROCESSOR|ODIN_MESSAGE_PROCESSOR"`
	RabbitMQConnectionData   *RabbitMQConnectionData   `json:"rabbitMqConnectionData"`
	ServiceBusConnectionData *ServiceBusConnectionData `json:"serviceBusConnectionData"`
}

func (s AMQPUpdate) Valid(validation *validation.Validation) {
	if s.Broker != connecting.AMQPBrokerRabbitMQ && s.Broker != connecting.AMQPBrokerServiceBus {
		_ = validation.SetError("broker", "should be one of [RABBIMQ,SERVICEBUS]")
	}

	if s.ProcessorType != message_processors.MessageProcessorTypeCDC && s.ProcessorType != message_processors.MessageProcessorTypeOdin {
		_ = validation.SetError("processorType", "should be one of [CDC_MESSAGE_PROCESSOR,ODIN_MESSAGE_PROCESSOR]")
	}

	if s.Broker == connecting.AMQPBrokerRabbitMQ && s.RabbitMQConnectionData == nil {
		_ = validation.SetError("rabbitMqConnectionData", "is required")
	}

	if s.Broker == connecting.AMQPBrokerServiceBus && s.ServiceBusConnectionData == nil {
		_ = validation.SetError("serviceBusConnectionData", "is required")
	}

	if s.RabbitMQConnectionData != nil && s.ServiceBusConnectionData != nil {
		_ = validation.SetError("*ConnectionData", "please data for only one broker type")
	}
}

func (s AMQPUpdate) ToMap() map[string]interface{} {
	var result interface{}
	if s.Broker == connecting.AMQPBrokerRabbitMQ {
		result = struct {
			ConnectionString string `model:"ConnectionString"`
			Broker           string `model:"Broker"`
			ProcessorType    string `model:"ProcessorType"`
			VirtualHost      string `model:"VirtualHost"`
			ExchangeName     string `model:"ExchangeName"`
			QueueName        string `model:"QueueName"`
			QueueKey         string `model:"QueueKey"`
		}{
			ConnectionString: s.ConnectionString,
			Broker:           s.Broker,
			ProcessorType:    s.ProcessorType,
			VirtualHost:      s.RabbitMQConnectionData.VirtualHost,
			ExchangeName:     s.RabbitMQConnectionData.ExchangeName,
			QueueName:        s.RabbitMQConnectionData.QueueName,
			QueueKey:         s.RabbitMQConnectionData.QueueKey,
		}
	} else {
		result = struct {
			ConnectionString string `model:"ConnectionString"`
			Broker           string `model:"Broker"`
			ProcessorType    string `model:"ProcessorType"`
			Topic            string `model:"Topic"`
			SubscriptionName string `model:"SubscriptionName"`
		}{
			ConnectionString: s.ConnectionString,
			Broker:           s.Broker,
			ProcessorType:    s.ProcessorType,
			Topic:            s.ServiceBusConnectionData.Topic,
			SubscriptionName: s.ServiceBusConnectionData.SubscriptionName,
		}
	}
	m, _ := model.Map(result)
	return m
}

type AMQPRequest struct {
	EventType string `json:"eventType" model:"event_type,omitempty"`
	Enabled   string `json:"enabled" model:"enabled,omitempty"`
	IssuerId  int    `json:"issuerId,string" model:"issuer_id,omitempty"`
}

func (request AMQPRequest) ToMap() map[string]interface{} {
	return helpers.ParseStructToMap(request)
}

type AMQPResponse struct {
	ID                       string                    `json:"id"`
	ConnectionString         string                    `json:"connectionString"`
	EventType                string                    `json:"eventType"`
	Enabled                  bool                      `json:"enabled"`
	CreatedAt                time.Time                 `json:"createdAt"`
	UpdatedAt                time.Time                 `json:"updatedAt"`
	DeletedAt                *time.Time                `json:"deletedAt,omitempty"`
	IssuerId                 int                       `json:"issuerId"`
	Broker                   string                    `json:"broker" example:"SERVICEBUS|RABBITMQ"`
	ProcessorType            string                    `json:"processorType" example:"CDC_MESSAGE_PROCESSOR|ODIN_MESSAGE_PROCESSOR"`
	RabbitMQConnectionData   *RabbitMQConnectionData   `json:"rabbitMqConnectionData,omitempty"`
	ServiceBusConnectionData *ServiceBusConnectionData `json:"serviceBusConnectionData,omitempty"`
}

func ConvertEntityToResponse(connection connecting.Connection) AMQPResponse {
	if connection.Broker == connecting.AMQPBrokerServiceBus {
		return AMQPResponse{
			ID:               connection.ID,
			ConnectionString: connection.GetMaskedConnectionString(),
			EventType:        connection.EventType,
			Enabled:          connection.Enabled,
			CreatedAt:        connection.CreatedAt,
			UpdatedAt:        connection.UpdatedAt,
			DeletedAt:        connection.DeletedAt,
			IssuerId:         connection.IssuerId,
			Broker:           connection.Broker,
			ProcessorType:    connection.ProcessorType,
			ServiceBusConnectionData: &ServiceBusConnectionData{
				Topic:            connection.Topic,
				SubscriptionName: connection.SubscriptionName,
			},
		}
	} else {
		return AMQPResponse{
			ID:               connection.ID,
			ConnectionString: connection.ConnectionString,
			EventType:        connection.EventType,
			Enabled:          connection.Enabled,
			CreatedAt:        connection.CreatedAt,
			UpdatedAt:        connection.UpdatedAt,
			DeletedAt:        connection.DeletedAt,
			IssuerId:         connection.IssuerId,
			Broker:           connection.Broker,
			ProcessorType:    connection.ProcessorType,
			RabbitMQConnectionData: &RabbitMQConnectionData{
				VirtualHost:  connection.VirtualHost,
				ExchangeName: connection.ExchangeName,
				QueueName:    connection.QueueName,
				QueueKey:     connection.QueueKey,
			},
		}
	}
}

type AMQPResponsePage struct {
	helpers.PageResponse
	Records []AMQPResponse `json:"records"`
}

func ConvertPagination(paginator helpers.Paginator) AMQPResponsePage {
	slice := make([]AMQPResponse, 0)
	records := paginator.Records.(*[]connecting.Connection)

	for _, rec := range *records {
		slice = append(slice, ConvertEntityToResponse(rec))
	}

	return AMQPResponsePage{
		PageResponse: helpers.PageResponse{
			TotalRecord: paginator.TotalRecord,
			TotalPage:   paginator.TotalPage,
			Offset:      paginator.Offset,
			Limit:       paginator.Limit,
			Page:        paginator.Page,
			PrevPage:    paginator.PrevPage,
			NextPage:    paginator.NextPage,
		},
		Records: slice,
	}
}

type ResendDeadLetterResponse struct {
	MessagesToBeSent int `json:"messagesToBeSent"`
}

type RetainedMessagesCount struct {
	RetainedMessages int32 `json:"retainedMessages"`
}

type ConnectionMetadata struct {
	LockDuration                              string                         `json:"lockDuration"`
	DefaultMessageTTL                         string                         `json:"defaultMessageTTL"`
	DeadLetteringOnMessageExpiration          bool                           `json:"deadLetteringOnMessageExpiration"`
	DeadLetteringOnFilterEvaluationExceptions bool                           `json:"deadLetteringOnFilterEvaluationExceptions"`
	MessageCount                              int64                          `json:"messageCount"`
	MaxDeliveryCount                          int32                          `json:"maxDeliveryCount"`
	Status                                    string                         `json:"status"`
	CreatedAt                                 string                         `json:"createdAt"`
	UpdatedAt                                 string                         `json:"updatedAt"`
	AccessedAt                                string                         `json:"accessedAt"`
	AutoDeleteOnIdle                          string                         `json:"autoDeleteOnIdle"`
	CountDetails                              ConnectionMetadataCountDetails `json:"countDetails"`
}

func (c *ConnectionMetadata) FromServiceBus(serviceBusEntity *servicebus.SubscriptionEntity) *ConnectionMetadata {
	if serviceBusEntity == nil {
		return c
	}

	c = &ConnectionMetadata{
		LockDuration:                              *serviceBusEntity.LockDuration,
		DefaultMessageTTL:                         *serviceBusEntity.DefaultMessageTimeToLive,
		DeadLetteringOnMessageExpiration:          *serviceBusEntity.DeadLetteringOnMessageExpiration,
		DeadLetteringOnFilterEvaluationExceptions: *serviceBusEntity.DeadLetteringOnFilterEvaluationExceptions,
		MessageCount:                              *serviceBusEntity.MessageCount,
		MaxDeliveryCount:                          *serviceBusEntity.MaxDeliveryCount,
		Status:                                    string(*serviceBusEntity.Status),
		CreatedAt:                                 serviceBusEntity.CreatedAt.String(),
		UpdatedAt:                                 serviceBusEntity.UpdatedAt.String(),
		AccessedAt:                                serviceBusEntity.AccessedAt.String(),
		AutoDeleteOnIdle:                          *serviceBusEntity.AutoDeleteOnIdle,
	}

	if serviceBusEntity.CountDetails != nil {
		c.CountDetails = ConnectionMetadataCountDetails{
			ActiveMessageCount:             *serviceBusEntity.CountDetails.ActiveMessageCount,
			DeadLetterMessageCount:         *serviceBusEntity.CountDetails.DeadLetterMessageCount,
			ScheduledMessageCount:          *serviceBusEntity.CountDetails.ScheduledMessageCount,
			TransferDeadLetterMessageCount: *serviceBusEntity.CountDetails.TransferDeadLetterMessageCount,
			TransferMessageCount:           *serviceBusEntity.CountDetails.TransferMessageCount,
		}
	}
	return c
}

type ConnectionMetadataCountDetails struct {
	ActiveMessageCount             int32 `json:"activeMessageCount"`
	DeadLetterMessageCount         int32 `json:"deadLetterMessageCount"`
	ScheduledMessageCount          int32 `json:"scheduledMessageCount"`
	TransferDeadLetterMessageCount int32 `json:"transferDeadLetterMessageCount"`
	TransferMessageCount           int32 `json:"transferMessageCount"`
}
