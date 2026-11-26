package event_controller

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Controller interface {
	UpdateEvents(tenantKey *tenants.TenantKey)
	UpdateEventByIssuerAndName(tenantKey *tenants.TenantKey, issuerId int, eventName string) error
}

func NewEventController(switcher Switcher, repository Repository) Controller {
	return controller{switcher: switcher, repository: repository}
}

type controller struct {
	switcher   Switcher
	repository Repository
}

func (c controller) UpdateEvents(tenantKey *tenants.TenantKey) {
	connections, err := c.repository.listAllConnections(tenantKey)
	if err != nil {
		logging.Logger.Error("[EVENT CONTROLLER] could not list enabled connections with CDC processor type", zap.Error(err))
		return
	}

	for _, connection := range connections {
		exists, err := c.existsWebHookByIssuerAndEvent(tenantKey, connection.IssuerId, connection.EventType)
		if err != nil {
			logging.Logger.Info("[EVENT CONTROLLER] could not find webhook", zap.Error(err))
			continue
		}

		c.switcher.set(connection.IssuerId, connection.EventType, exists)
	}
}

func (c controller) UpdateEventByIssuerAndName(tenantKey *tenants.TenantKey, issuerId int, eventName string) error {
	err := c.repository.existsConnectionByIssuerAndEventName(tenantKey, issuerId, eventName)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.switcher.Off(issuerId, eventName)
			return nil
		} else {
			return err
		}
	}

	exists, err := c.existsWebHookByIssuerAndEvent(tenantKey, issuerId, eventName)
	if err != nil {
		return err
	}
	c.switcher.set(issuerId, eventName, exists)
	return nil
}

func (c controller) existsWebHookByIssuerAndEvent(tenantKey *tenants.TenantKey, issuerId int, eventName string) (bool, error) {
	err := c.repository.existsWebHookByIssuerAndEventName(tenantKey, issuerId, eventName)
	if err == nil {
		return true, nil
	} else {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		} else {
			return false, err
		}
	}
}
