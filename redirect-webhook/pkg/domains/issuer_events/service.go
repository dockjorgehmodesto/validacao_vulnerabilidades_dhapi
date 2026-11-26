package issuer_events

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Create(tenantKey *tenants.TenantKey, vault IssuerEvent) IssuerEvent
	Update(tenantKey *tenants.TenantKey, id string, updateValues map[string]interface{}) IssuerEvent
	FindById(tenantKey *tenants.TenantKey, id string) IssuerEvent
	FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, queryExample map[string]interface{}) helpers.Paginator
	Delete(tenantKey *tenants.TenantKey, id string)
	FindByEventNameAndIssuerId(key *tenants.TenantKey, eventName string, issuerId int) IssuerEvent
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return serviceImpl{repository: repository}
}

func (s serviceImpl) Create(tenantKey *tenants.TenantKey, issuerEvent IssuerEvent) IssuerEvent {
	if s.hasDuplicateByNameAndIssuerId(tenantKey, issuerEvent.Name, issuerEvent.IssuerId) {
		Error{}.HandleNameAlreadyTaken()
	}
	savedIssuerEvent, e := s.repository.Save(tenantKey, issuerEvent)
	Error{}.HandleDatabaseError(e)
	return savedIssuerEvent
}

func (s serviceImpl) Update(tenantKey *tenants.TenantKey, id string, updateValues map[string]interface{}) IssuerEvent {
	s.FindById(tenantKey, id)

	connectionId := updateValues["connection_id"].(string)
	name := updateValues["name"].(string)
	issuerId := updateValues["issuer_id"].(int)
	if !s.connectionIdIsValid(tenantKey, connectionId, name) {
		Error{}.HandleInvalidConnection()
	}

	if s.existsByConnectionAndIssuerIdIgnoringId(tenantKey, connectionId, issuerId, &id) {
		Error{}.HandleConnectionAlreadyAssociatedToAnEvent()
	}

	hasDuplicate := s.hasDuplicateByNameIgnoringId(tenantKey, name, id, issuerId)
	if hasDuplicate {
		Error{}.HandleNameAlreadyTaken()
	}

	e := s.repository.PartialUpdate(tenantKey, id, updateValues)
	Error{}.HandleDatabaseError(e)
	return s.FindById(tenantKey, id)
}

func (s serviceImpl) FindById(tenantKey *tenants.TenantKey, id string) IssuerEvent {
	issuerEvent, e := s.repository.FindById(tenantKey, id)
	Error{}.HandleSearch(e)
	return issuerEvent
}

func (s serviceImpl) FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, queryExample map[string]interface{}) helpers.Paginator {
	paginator, e := s.repository.FindAllWithPagination(tenantKey, page, limit, sort, queryExample)
	Error{}.HandlePagination(paginator, e)
	return paginator
}

func (s serviceImpl) hasDuplicateByNameAndIssuerId(tenantKey *tenants.TenantKey, issuerEventName string, issuerId int) bool {
	_, err := s.repository.FindByNameAndIssuerId(tenantKey, issuerEventName, issuerId)
	if err == nil {
		return true
	}
	if !gorm.IsRecordNotFoundError(err) {
		Error{}.HandleDatabaseErrorWithMessage(err, "error finding issuer event by name")
	}
	return false
}

func (s serviceImpl) hasDuplicateByNameIgnoringId(tenantKey *tenants.TenantKey, issuerEventName string, ignoreId string, issuerId int) bool {
	issuerEvent, err := s.repository.FindByNameAndIssuerId(tenantKey, issuerEventName, issuerId)
	if err == nil {
		return issuerEvent.Id != ignoreId
	}
	if !gorm.IsRecordNotFoundError(err) {
		Error{}.HandleDatabaseErrorWithMessage(err, "error finding issuer event by name")
	}
	return false
}

func (s serviceImpl) Delete(tenantKey *tenants.TenantKey, id string) {
	s.FindById(tenantKey, id)
	err := s.repository.Delete(tenantKey, id)
	Error{}.HandleDatabaseErrorWithMessage(err, "could not delete issuer event")
}

func (s serviceImpl) connectionIdIsValid(tenantKey *tenants.TenantKey, connectionId string, eventName string) bool {
	err := s.repository.validateConnectionByIdAndEventName(tenantKey, connectionId, eventName)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		} else {
			Error{}.HandleDatabaseErrorWithMessage(err, "could not find connection with given id and event type")
		}
	}
	return true
}

func (s serviceImpl) existsByConnectionAndIssuerIdIgnoringId(tenantKey *tenants.TenantKey, connectionId string, issuerId int, issuerEventId *string) bool {
	err := s.repository.existsByConnectionAndIssuerIdIgnoringId(tenantKey, connectionId, issuerId, issuerEventId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		} else {
			Error{}.HandleDatabaseErrorWithMessage(err, "could not validate existing issuer event")
		}
	}
	return true
}

func (s serviceImpl) FindByEventNameAndIssuerId(key *tenants.TenantKey, eventName string, issuerId int) IssuerEvent {
	issuerEvent, err := s.repository.FindByNameAndIssuerId(key, eventName, issuerId)
	if err != nil {
		Error{}.HandleEventNotSetupForIssuer()
	}
	return issuerEvent
}
