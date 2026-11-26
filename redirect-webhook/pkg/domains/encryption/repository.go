package encryption

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"context"
	"github.com/jinzhu/gorm"
	gorm2 "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"time"
)

type Repository interface {
	FindById(key *tenants.TenantKey, id string) (KeyVault, error)
	FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error)
	Save(key *tenants.TenantKey, keyvault KeyVault) (KeyVault, error)
	Update(key *tenants.TenantKey, keyvault KeyVault) (KeyVault, error)
	FindFirst(key *tenants.TenantKey) (KeyVault, error)
	FindMostRecentAndActive(key *tenants.TenantKey) (KeyVault, error)
	FindMostRecentAndActiveByIssuerId(ctx context.Context, key *tenants.TenantKey, issuerId int) (KeyVault, error)
}

type RepositoryImpl struct{}

func NewRepository() RepositoryImpl {
	return RepositoryImpl{}
}

func (r RepositoryImpl) FindById(key *tenants.TenantKey, id string) (KeyVault, error) {
	out := KeyVault{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.First(&out, "id = ? and issuer_id = ?", id, key.IssuerId).Error
	})
	return out, e
}

func (r RepositoryImpl) FindFirst(key *tenants.TenantKey) (KeyVault, error) {
	out := KeyVault{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Where("issuer_id = ?", key.IssuerId).First(&out).Error
	})
	return out, e
}

func (r RepositoryImpl) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error) {
	var paginator helpers.Paginator
	out := make([]KeyVault, 0)

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

func (r RepositoryImpl) Save(tenantKey *tenants.TenantKey, keyvault KeyVault) (KeyVault, error) {
	if keyvault.IssuerId == 0 {
		keyvault.IssuerId = tenantKey.IssuerId
	}
	e := tenants.GetConnection(tenantKey, func(db *gorm.DB) (ex error) {
		return db.Save(&keyvault).Error
	})
	return keyvault, e
}

func (r RepositoryImpl) Update(key *tenants.TenantKey, keyvault KeyVault) (KeyVault, error) {
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Where("id = ? and issuer_id = ?", keyvault.ID, key.IssuerId).Save(&keyvault).Error
	})
	return keyvault, err
}

func (r RepositoryImpl) FindMostRecentAndActive(key *tenants.TenantKey) (KeyVault, error) {
	out := KeyVault{}
	now := time.Now()
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Where("enabled = true and issuer_key_expiration_date > ? and aes_key_expiration_date > ? and issuer_id = ?", now, now, key.IssuerId).Order("issuer_key_expiration_date DESC").First(&out).Error
	})
	return out, e
}

func (r RepositoryImpl) FindMostRecentAndActiveByIssuerId(ctx context.Context, key *tenants.TenantKey, issuerId int) (KeyVault, error) {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "findKeyvaultOnGorm", tracer.Tag("repository", "gorm"))
	defer ddTrace.Finish()
	out := KeyVault{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		db = gorm2.WithContext(ctx, db)
		return db.Where("enabled = true and issuer_id = ?", issuerId).Order("issuer_key_expiration_date DESC").First(&out).Error
	})
	return out, e
}
