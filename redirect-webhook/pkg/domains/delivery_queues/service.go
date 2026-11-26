package delivery_queues

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Service interface {
	Create(tenantKey *tenants.TenantKey, deliveryQueue DeliveryQueue) DeliveryQueue
	Update(tenantKey *tenants.TenantKey, id string, updateValues map[string]interface{}) DeliveryQueue
	FindById(tenantKey *tenants.TenantKey, id string) DeliveryQueue
	FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, queryExample map[string]interface{}) helpers.Paginator
	Delete(tenantKey *tenants.TenantKey, id string)
	FindDeliveryQueue(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error)
	CreateList(tenantKey *tenants.TenantKey, deliveryQueueList DeliveryQueueList) []DeliveryQueue
	DeleteAll(tenantKey *tenants.TenantKey, id string)
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return serviceImpl{repository: repository}
}

func (s serviceImpl) Create(tenantKey *tenants.TenantKey, deliveryQueue DeliveryQueue) DeliveryQueue {
	if s.hasDuplicateByEventTypeAndIssuerId(tenantKey, deliveryQueue.EventType, deliveryQueue.IssuerId) {
		Error{}.HandleEventTypeAlreadyConfigured()
	}
	savedDeliveryQueue, e := s.repository.Save(tenantKey, deliveryQueue)
	Error{}.HandleDatabaseError(e)
	return savedDeliveryQueue
}

func (s serviceImpl) Update(tenantKey *tenants.TenantKey, id string, updateValues map[string]interface{}) DeliveryQueue {
	s.FindById(tenantKey, id)
	issuerId := updateValues["issuer_id"].(int)
	eventType := updateValues["event_type"].(string)
	hasDuplicate := s.hasDuplicateByEventTypeIgnoringId(tenantKey, eventType, id, issuerId)
	if hasDuplicate {
		Error{}.HandleEventTypeAlreadyConfigured()
	}

	e := s.repository.PartialUpdate(tenantKey, id, updateValues)
	Error{}.HandleDatabaseError(e)
	return s.FindById(tenantKey, id)
}

func (s serviceImpl) FindById(tenantKey *tenants.TenantKey, id string) DeliveryQueue {
	deliveryQueue, e := s.repository.FindById(tenantKey, id)
	Error{}.HandleSearch(e)
	return deliveryQueue
}

func (s serviceImpl) FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, queryExample map[string]interface{}) helpers.Paginator {
	paginator, e := s.repository.FindAllWithPagination(tenantKey, page, limit, sort, queryExample)
	Error{}.HandlePagination(paginator, e)
	return paginator
}

func (s serviceImpl) hasDuplicateByEventTypeAndIssuerId(tenantKey *tenants.TenantKey, deliveryQueueEventType string, issuerId int) bool {
	_, err := s.repository.FindByEventTypeAndIssuerId(tenantKey, deliveryQueueEventType, issuerId)
	if err == nil {
		return true
	}
	if !gorm.IsRecordNotFoundError(err) {
		Error{}.HandleDatabaseErrorWithMessage(err, "error finding delivery queue by event type")
	}
	return false
}

func (s serviceImpl) hasDuplicateByEventTypeIgnoringId(tenantKey *tenants.TenantKey, eventType string, ignoreId string, issuerId int) bool {
	deliveryQueue, err := s.repository.FindByEventTypeAndIssuerId(tenantKey, eventType, issuerId)
	if err == nil {
		return deliveryQueue.Id != ignoreId
	}
	if !gorm.IsRecordNotFoundError(err) {
		Error{}.HandleDatabaseErrorWithMessage(err, "error finding delivery queue by event type")
	}
	return false
}

func (s serviceImpl) Delete(tenantKey *tenants.TenantKey, id string) {
	queue := s.FindById(tenantKey, id)
	err := s.repository.Delete(tenantKey, id, queue.EventType)
	Error{}.HandleDatabaseErrorWithMessage(err, "could not delete delivery queue")
}

func (s serviceImpl) FindDeliveryQueue(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error) {
	deliveryQueue, err := s.repository.FindByEventTypeAndIssuerIdAndEnabled(tenantKey, eventType, issuerId)
	if err != nil {
		return deliveryQueue, err
	}
	return deliveryQueue, nil
}

func (s serviceImpl) CreateList(tenantKey *tenants.TenantKey, deliveryQueue DeliveryQueueList) []DeliveryQueue {
	var deliveryQueueList []DeliveryQueue
	if deliveryQueue.EventTypeList == nil {
		deliveryQueue.EventTypeList = EventTypeList
	}

	for _, eventType := range deliveryQueue.EventTypeList {
		q := DeliveryQueue{
			IssuerId: deliveryQueue.IssuerId,
			EventType: eventType,
			ExchangeName: deliveryQueue.ExchangeName,
			QueueName:    deliveryQueue.QueueName,
			QueueKey:     deliveryQueue.QueueKey,
			Enabled: deliveryQueue.Enabled,
		}
		deliveryQueueList = append(deliveryQueueList, q)
	}
	savedDeliveryQueueList, err := s.repository.DeleteIfExistsAndSaveAll(tenantKey, deliveryQueueList)
	Error{}.HandleDatabaseErrorWithMessage(err, "error creating delivery queue list")
	return savedDeliveryQueueList
}

func (s serviceImpl) DeleteAll(tenantKey *tenants.TenantKey, issuerId string) {
	intIssuerId, err := strconv.Atoi(issuerId)
	if err != nil {
		Error{}.HandleDatabaseErrorWithMessage(err, "invalid id")
	}
	err = s.repository.DeleteAll(tenantKey, intIssuerId)
	Error{}.HandleDatabaseErrorWithMessage(err, "could not delete delivery queue")
}
