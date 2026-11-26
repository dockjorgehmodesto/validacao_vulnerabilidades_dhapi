package delivery_queues

import (
	"conductor.com.br/pier/hooks/pkg/helpers"
)

type DeliveryQueueRequest struct {
	IssuerId     *int    `json:"issuerId,string" model:"issuer_id,omitempty"`
	EventType    *string `json:"eventType" model:"event_type,omitempty"`
	ExchangeName *string `json:"exchangeName" model:"exchange_name,omitempty"`
	QueueName    *string `json:"queueName" model:"queue_name,omitempty"`
	QueueKey     *string `json:"queueKey" model:"queue_key,omitempty"`
	Enabled      *string `json:"enabled" model:"enabled,omitempty"`
}

func (r DeliveryQueueRequest) ToMap() map[string]interface{} {
	return helpers.ParseStructToMap(r)
}

type DeliveryQueueResponse struct {
	Id           string `json:"id"`
	IssuerId     int    `json:"issuerId"`
	EventType    string `json:"eventType"`
	ExchangeName string `json:"exchangeName"`
	QueueName    string `json:"queueName"`
	QueueKey     string `json:"queueKey"`
	Enabled      bool   `json:"enabled"`
}

func ConvertEntityToResponse(event DeliveryQueue) DeliveryQueueResponse {
	return DeliveryQueueResponse{
		Id:           event.Id,
		IssuerId:     event.IssuerId,
		EventType:    event.EventType,
		ExchangeName: event.ExchangeName,
		QueueName:    event.QueueName,
		QueueKey:     event.QueueKey,
		Enabled:      event.Enabled,
	}
}

func ConvertListToResponse(deliveryQueueList []DeliveryQueue) []DeliveryQueueResponse {
	var response []DeliveryQueueResponse
	for _, deliveryQueue:= range deliveryQueueList {
		response = append(response, DeliveryQueueResponse{
			Id:           deliveryQueue.Id,
			IssuerId:     deliveryQueue.IssuerId,
			EventType:    deliveryQueue.EventType,
			ExchangeName: deliveryQueue.ExchangeName,
			QueueName:    deliveryQueue.QueueName,
			QueueKey:     deliveryQueue.QueueKey,
			Enabled:      deliveryQueue.Enabled,
		})
	}
	return response
}

type DeliveryQueueResponsePage struct {
	helpers.PageResponse
	Records []DeliveryQueueResponse `json:"records"`
}

func ConvertPagination(paginator helpers.Paginator) DeliveryQueueResponsePage {
	slice := make([]DeliveryQueueResponse, 0)
	records := paginator.Records.(*[]DeliveryQueue)

	for _, rec := range *records {
		slice = append(slice, ConvertEntityToResponse(rec))
	}

	return DeliveryQueueResponsePage{
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

type DeliveryQueueNameResponse struct {
	ExchangeName string `json:"exchangeName"`
}

type DeliveryQueueNameResponsePage struct {
	helpers.PageResponse
	Records []DeliveryQueueNameResponse `json:"records"`
}

type DeliveryQueuePersist struct {
	IssuerId     int    `json:"issuerId" valid:"Required"`
	EventType    string `json:"eventType" valid:"Required;MaxSize(250)"`
	ExchangeName string `json:"exchangeName" valid:"Required;MaxSize(250)"`
	QueueName    string `json:"queueName" model:"Required;MaxSize(250)"`
	QueueKey     string `json:"queueKey" model:"Required;MaxSize(250)"`
	Enabled      bool   `json:"enabled" valid:"Required"`
}

type DeliveryQueueList struct {
	IssuerId      int      `json:"issuerId" valid:"Required"`
	EventTypeList []string `json:"eventTypeList"`
	ExchangeName  string   `json:"exchangeName" valid:"Required;MaxSize(250)"`
	QueueName     string   `json:"queueName" model:"Required;MaxSize(250)"`
	QueueKey      string   `json:"queueKey" model:"Required;MaxSize(250)"`
	Enabled       bool     `json:"enabled" valid:"Required"`
}

func (i DeliveryQueuePersist) ToEntity() DeliveryQueue {
	return DeliveryQueue{
		IssuerId:     i.IssuerId,
		EventType:    i.EventType,
		ExchangeName: i.ExchangeName,
		QueueName:    i.QueueName,
		QueueKey:     i.QueueKey,
		Enabled:      i.Enabled,
	}
}

type DeliveryQueueUpdate struct {
	IssuerId     int    `json:"issuerId" valid:"Required" model:"issuer_id"`
	EventType    string `json:"eventType" valid:"Required;MaxSize(250)" model:"event_type"`
	ExchangeName string `json:"exchangeName" valid:"Required;MaxSize(250)" model:"exchange_name"`
	QueueName    string `json:"queueName" valid:"Required;MaxSize(250)" model:"queue_name"`
	QueueKey     string `json:"queueKey" valid:"Required;MaxSize(250)" model:"queue_key"`
	Enabled      bool   `json:"enabled" valid:"Required;" model:"enabled"`
}

func (i DeliveryQueueUpdate) ToMap() map[string]interface{} {
	return helpers.ParseStructToMap(i)
}
