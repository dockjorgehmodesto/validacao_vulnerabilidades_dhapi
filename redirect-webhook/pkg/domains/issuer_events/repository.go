package issuer_events

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	FindById(tenantKey *tenants.TenantKey, id string) (IssuerEvent, error)
	FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error)
	Save(tenantKey *tenants.TenantKey, issuerEvent IssuerEvent) (IssuerEvent, error)
	PartialUpdate(tenantKey *tenants.TenantKey, issuerEventId string, updateValues map[string]interface{}) error
	FindByNameAndIssuerId(tenantKey *tenants.TenantKey, issuerEventName string, issuerId int) (IssuerEvent, error)
	Delete(tenantKey *tenants.TenantKey, id string) error
	validateConnectionByIdAndEventName(tenantKey *tenants.TenantKey, connectionId string, eventName string) error
	existsByConnectionAndIssuerIdIgnoringId(tenantKey *tenants.TenantKey, connectionId string, issuerId int, issuerEventId *string) error
}

type repositoryImpl struct{}

func NewRepository() Repository {
	return repositoryImpl{}
}

func (r repositoryImpl) FindById(tenantKey *tenants.TenantKey, id string) (IssuerEvent, error) {
	out := IssuerEvent{}
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.First(&out, "id = ?", id).Error
	})
	return out, e
}

func (r repositoryImpl) FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error) {
	var paginator helpers.Paginator
	out := make([]IssuerEvent, 0)

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

func (r repositoryImpl) Save(tenantKey *tenants.TenantKey, issuerEvent IssuerEvent) (IssuerEvent, error) {
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Save(&issuerEvent).Error
	})
	return issuerEvent, e
}

func (r repositoryImpl) PartialUpdate(tenantKey *tenants.TenantKey, issuerEventId string, updateValues map[string]interface{}) error {
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Model(IssuerEvent{}).Where("id = ?", issuerEventId).Update(updateValues).Error
	})
}

func (r repositoryImpl) FindByNameAndIssuerId(tenantKey *tenants.TenantKey, issuerEventName string, issuerId int) (IssuerEvent, error) {
	out := IssuerEvent{}
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.First(&out, "name = ? and issuer_id = ?", issuerEventName, issuerId).Error
	})
	return out, e
}

func (r repositoryImpl) Delete(tenantKey *tenants.TenantKey, id string) error {
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Where("id = ?", id).Delete(IssuerEvent{}).Error
	})
}

func (r repositoryImpl) validateConnectionByIdAndEventName(tenantKey *tenants.TenantKey, connectionId string, eventName string) error {
	type result struct {
		Id string
	}
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		id := result{}
		return db.Table("amqp_connections").Select("id").Where("id = ? and event_type = ?", connectionId, eventName).Scan(&id).Error
	})
}

func (r repositoryImpl) existsByConnectionAndIssuerIdIgnoringId(tenantKey *tenants.TenantKey, connectionId string, issuerId int, issuerEventId *string) error {
	type result struct {
		Id string
	}
	return tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		id := result{}
		appendStatement := "id is not null"
		if issuerEventId != nil {
			appendStatement = fmt.Sprintf(`id != '%s'`, *issuerEventId)
		}
		return db.Table("issuer_events").Select("id").Where("connection_id = ? and issuer_id = ? and "+appendStatement, connectionId, issuerId).Scan(&id).Error
	})
}
