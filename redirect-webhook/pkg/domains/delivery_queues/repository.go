package delivery_queues

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	FindById(tenantKey *tenants.TenantKey, id string) (DeliveryQueue, error)
	FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error)
	Save(tenantKey *tenants.TenantKey, deliveryQueue DeliveryQueue) (DeliveryQueue, error)
	PartialUpdate(tenantKey *tenants.TenantKey, deliveryQueueId string, updateValues map[string]interface{}) error
	FindByEventTypeAndIssuerId(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error)
	Delete(tenantKey *tenants.TenantKey, id string, eventType string) error
	FindByEventTypeAndIssuerIdAndEnabled(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error)
	DeleteIfExistsAndSaveAll(tenantKey *tenants.TenantKey, deliveryQueueList []DeliveryQueue) ([]DeliveryQueue, error)
	DeleteAll(tenantKey *tenants.TenantKey, issuerId int) error
}

type repositoryImpl struct{}

func NewRepository() Repository {
	return repositoryImpl{}
}

func (r repositoryImpl) FindById(tenantKey *tenants.TenantKey, id string) (DeliveryQueue, error) {
	out := DeliveryQueue{}
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.First(&out, "id = ?", id).Error
	})
	return out, e
}

func (r repositoryImpl) FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error) {
	var paginator helpers.Paginator
	out := make([]DeliveryQueue, 0)

	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		dbWithExample := db
		if example != nil {
			dbWithExample = db.Where(example)
		}

		paginator = *helpers.Paging(&helpers.Param{
			DB:      dbWithExample,
			Page:    page,
			Limit:   limit,
			OrderBy: sort,
			ShowSQL: global.ShowSQL,
		}, &out)
		return nil
	})

	return paginator, e
}

func (r repositoryImpl) Save(tenantKey *tenants.TenantKey, deliveryQueue DeliveryQueue) (DeliveryQueue, error) {
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Save(&deliveryQueue).Error
	})
	return deliveryQueue, e
}

func (r repositoryImpl) PartialUpdate(tenantKey *tenants.TenantKey, deliveryQueueId string, updateValues map[string]interface{}) error {
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Model(DeliveryQueue{}).Where("id = ?", deliveryQueueId).Update(updateValues).Error
	})
}

func (r repositoryImpl) FindByEventTypeAndIssuerId(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error) {
	out := DeliveryQueue{}
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.First(&out, "event_type = ? and issuer_id = ?", eventType, issuerId).Error
	})
	return out, e
}

func (r repositoryImpl) Delete(tenantKey *tenants.TenantKey, id string, eventType string) error {
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Unscoped().Where("id = ?", id).Delete(DeliveryQueue{}).Error
	})
}

func (r repositoryImpl) DeleteAll(tenantKey *tenants.TenantKey, issuerId int) error {
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Unscoped().Where("issuer_id = ?", issuerId).Delete(DeliveryQueue{}).Error
	})
}

func (r repositoryImpl) DeleteIfExistsAndSaveAll(tenantKey *tenants.TenantKey, deliveryQueueList []DeliveryQueue) ([]DeliveryQueue, error) {
	var insertedDeliveryQueues []DeliveryQueue
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		for _, deliveryQueue := range deliveryQueueList {
			err := db.Unscoped().Where("issuer_id = ? and event_type = ?", tenantKey.IssuerId, deliveryQueue.EventType).Delete(DeliveryQueue{}).Error
			if err != nil {
				return err
			}
			err = db.Save(&deliveryQueue).Error
			if err != nil {
				return err
			}
			insertedDeliveryQueues = append(insertedDeliveryQueues, deliveryQueue)
		}
		return nil
	})
	return insertedDeliveryQueues, e
}

func (r repositoryImpl) FindByEventTypeAndIssuerIdAndEnabled(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error) {
	out := DeliveryQueue{}
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.First(&out, "event_type = ? and issuer_id = ? and enabled = true", eventType, issuerId).Error
	})
	return out, e
}
