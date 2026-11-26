package amqp

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	FindById(key *tenants.TenantKey, id string) (connecting.Connection, error)
	FindByEventType(key *tenants.TenantKey, eventType string) (connecting.Connection, error)
	FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error)
	Save(key *tenants.TenantKey, connection connecting.Connection) (connecting.Connection, error)
	Update(key *tenants.TenantKey, connection connecting.Connection) (connecting.Connection, error)
	UpdateStatus(key *tenants.TenantKey, connection connecting.Connection) error
	PartialUpdate(key *tenants.TenantKey, id string, values map[string]interface{}) (connecting.Connection, error)
	ListEnabled(key *tenants.TenantKey) ([]connecting.Connection, error)
	ListEnabledByIssuer(key *tenants.TenantKey) ([]connecting.Connection, error)
	Delete(key *tenants.TenantKey, id string) error
	DeleteByIssuer(key *tenants.TenantKey) error
}

type RepositoryImpl struct{}

func NewRepository() RepositoryImpl {
	return RepositoryImpl{}
}

func (r RepositoryImpl) FindById(key *tenants.TenantKey, id string) (connecting.Connection, error) {
	out := connecting.Connection{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.First(&out, "id = ?", id).Error
	})
	return out, e
}

func (r RepositoryImpl) FindByEventType(key *tenants.TenantKey, eventType string) (connecting.Connection, error) {
	out := connecting.Connection{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.First(&out, "event_type = ? and issuer_id = ?", eventType, key.IssuerId).Error
	})
	return out, e
}

func (r RepositoryImpl) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error) {
	var paginator helpers.Paginator
	out := make([]connecting.Connection, 0)

	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		dbWithExample := db
		if example != nil {
			dbWithExample = db.Where(example)
		}
		dbWithExample.Where("issuer_id = ?", key.IssuerId)

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

func (r RepositoryImpl) Save(key *tenants.TenantKey, connection connecting.Connection) (connecting.Connection, error) {
	connection.IssuerId = key.IssuerId
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Save(&connection).Error
	})
	return connection, e
}

func (r RepositoryImpl) Update(key *tenants.TenantKey, connection connecting.Connection) (connecting.Connection, error) {
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Model(&connection).Update(&connection).Error
	})
	return connection, err
}

func (r RepositoryImpl) UpdateStatus(key *tenants.TenantKey, connection connecting.Connection) error {
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Model(&connection).Exec("UPDATE AMQP_CONNECTIONS SET ENABLED = ? WHERE ID = ? AND DELETED_AT IS NULL", connection.Enabled, connection.ID).Error
	})
	return err
}

func (r RepositoryImpl) PartialUpdate(key *tenants.TenantKey, id string, values map[string]interface{}) (connecting.Connection, error) {
	connection := connecting.Connection{}
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Model(&connection).Where("id = ?", id).Update(values).Error
	})
	return connection, err
}

func (r RepositoryImpl) ListEnabled(key *tenants.TenantKey) ([]connecting.Connection, error) {
	out := make([]connecting.Connection, 0)
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Where("enabled = ?", true).Find(&out).Error
	})
	return out, e
}

func (r RepositoryImpl) ListEnabledByIssuer(key *tenants.TenantKey) ([]connecting.Connection, error) {
	out := make([]connecting.Connection, 0)
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Where("enabled = ? and issuer_id = ?", true, key.IssuerId).Find(&out).Error
	})
	return out, e
}

func (r RepositoryImpl) Delete(key *tenants.TenantKey, id string) error {
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Where("id = ?", id).Delete(connecting.Connection{}).Error
	})
	return err
}

func (r RepositoryImpl) DeleteByIssuer(key *tenants.TenantKey) error {
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Where("issuer_id  = ?", key.IssuerId).Delete(connecting.Connection{}).Error
	})
	return err
}
