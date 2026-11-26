package event_controller

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"github.com/jinzhu/gorm"
)

type simplifiedConnection struct {
	Id        string
	EventType string
	IssuerId  int
}

type Repository interface {
	listAllConnections(tenantKey *tenants.TenantKey) ([]simplifiedConnection, error)
	existsWebHookByIssuerAndEventName(tenantKey *tenants.TenantKey, issuerId int, eventType string) error
	existsConnectionByIssuerAndEventName(tenantKey *tenants.TenantKey, issuerId int, eventName string) error
}

func NewRepository() Repository {
	return repository{}
}

type repository struct{}

func (r repository) listAllConnections(tenantKey *tenants.TenantKey) ([]simplifiedConnection, error) {
	result := make([]simplifiedConnection, 0)
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Table("amqp_connections").
			Select("id, event_type, issuer_id").
			Where("enabled = 'true' and processor_type='CDC_MESSAGE_PROCESSOR' and deleted_at is null").
			Scan(&result).Error
	})
	return result, e
}

func (r repository) existsWebHookByIssuerAndEventName(tenantKey *tenants.TenantKey, issuerId int, eventType string) error {
	webhookId := struct{ Id string }{}
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Table("web_hooks").
			Select("id").
			Where("deleted_at is null and event_type = ? and issuer_id = ?", eventType, issuerId).
			Scan(&webhookId).Error
	})
}

func (r repository) existsConnectionByIssuerAndEventName(tenantKey *tenants.TenantKey, issuerId int, eventName string) error {
	connectionId := struct{ Id string }{}
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Table("amqp_connections").
			Select("id").
			Where("enabled = 'true' and processor_type='CDC_MESSAGE_PROCESSOR' and deleted_at is null and event_type = ? and issuer_id = ?", eventName, issuerId).
			Scan(&connectionId).Error
	})
}
