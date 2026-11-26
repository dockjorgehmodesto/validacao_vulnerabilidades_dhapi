package webhooks

import (
	"context"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/jinzhu/gorm"
	gorm2 "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Repository interface {
	FindByIdWithPreload(key *tenants.TenantKey, id string) (WebHook, error)
	FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example *WebHook) (helpers.Paginator, error)
	Save(key *tenants.TenantKey, webHook WebHook) (WebHook, error)
	Update(key *tenants.TenantKey, id string, updateMap map[string]interface{}) error
	FindByType(key *tenants.TenantKey, eventType string) (WebHook, error)
	FindByTypeWithPreload(key *tenants.TenantKey, eventType string) (WebHook, error)
	Delete(tenant *tenants.TenantKey, id string, eventType string) error
	FindAllByTypeAndIssuerId(ctx context.Context, key *tenants.TenantKey, eventType string, id int) ([]WebHook, error)
	FindByTypeAndUrl(context context.Context, tenantKey *tenants.TenantKey, eventType string, url string) (WebHook, error)
}

type RepositoryImpl struct{}

func NewRepository() RepositoryImpl {
	return RepositoryImpl{}
}

func (r RepositoryImpl) FindByIdWithPreload(key *tenants.TenantKey, id string) (WebHook, error) {
	out := WebHook{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Set("gorm:auto_preload", true).First(&out, "id = ? and issuer_id = ?", id, key.IssuerId).Error
	})
	return out, e
}

func (r RepositoryImpl) FindByType(key *tenants.TenantKey, eventType string) (WebHook, error) {
	out := WebHook{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.First(&out, "event_type = ? and issuer_id = ?", eventType, key.IssuerId).Error
	})
	return out, e
}

func (r RepositoryImpl) FindByTypeWithPreload(key *tenants.TenantKey, eventType string) (WebHook, error) {
	out := WebHook{}
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Set("gorm:auto_preload", true).First(&out, "event_type = ? and issuer_id = ?", eventType, key.IssuerId).Error
	})
	return out, e
}

func (r RepositoryImpl) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example *WebHook) (helpers.Paginator, error) {
	var paginator helpers.Paginator
	out := make([]WebHook, 0)

	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		dbWithExample := db
		if example == nil {
			example = &WebHook{IssuerId: key.IssuerId}
		} else {
			example.IssuerId = key.IssuerId
		}
		dbWithExample = db.Where(example)

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

func (r RepositoryImpl) Save(key *tenants.TenantKey, webHook WebHook) (WebHook, error) {
	webHook.IssuerId = key.IssuerId
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Save(&webHook).Error
	})
	return webHook, e
}

func (r RepositoryImpl) Update(key *tenants.TenantKey, id string, updateMap map[string]interface{}) error {
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Model(WebHook{}).Where("id = ? and issuer_id = ?", id, key.IssuerId).Updates(updateMap).Error
	})
	return err
}

func (r RepositoryImpl) Delete(tenant *tenants.TenantKey, id string, _ string) error {
	err := tenants.GetConnection(tenant, func(db *gorm.DB) (ex error) {
		return db.Where("id = ?", id).Delete(WebHook{}).Error
	})
	return err
}

func (r RepositoryImpl) FindAllByTypeAndIssuerId(ctx context.Context, key *tenants.TenantKey, eventType string, issuerId int) ([]WebHook, error) {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "findWebhooksOnGorm", tracer.Tag("repository", "gorm"))
	defer ddTrace.Finish()
	result := make([]WebHook, 0)
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		db = gorm2.WithContext(ctx, db)
		return db.Set("gorm:auto_preload", true).Where("event_type = ? and issuer_id = ?", eventType, issuerId).Find(&result).Error
	})
	return result, err
}

func (r RepositoryImpl) FindByTypeAndUrl(ctx context.Context, key *tenants.TenantKey, eventType string, url string) (WebHook, error) {
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "findWebhooksOnGorm", tracer.Tag("repository", "gorm"))
	defer ddTrace.Finish()
	result := WebHook{}
	err := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		db = gorm2.WithContext(ctx, db)
		return db.Set("gorm:auto_preload", true).Where("event_type = ? and url = ? and issuer_id = ?", eventType, url, key.IssuerId).Find(&result).Error
	})
	return result, err
}

type RepositoryHeader interface {
	Delete(tenant *tenants.TenantKey, id string) error
	DeleteByWebhookId(tenant *tenants.TenantKey, webhookId string) error
	Save(key *tenants.TenantKey, headers Header) (Header, error)
}

type RepositoryHeaderImpl struct{}

func NewRepositoryHeader() RepositoryHeaderImpl {
	return RepositoryHeaderImpl{}
}

func (r RepositoryHeaderImpl) Delete(tenant *tenants.TenantKey, id string) error {
	err := tenants.GetConnection(tenant, func(db *gorm.DB) (ex error) {
		return db.Where("id = ?", id).Delete(Header{}).Error
	})
	return err
}

func (r RepositoryHeaderImpl) DeleteByWebhookId(tenant *tenants.TenantKey, webhookId string) error {
	err := tenants.GetConnection(tenant, func(db *gorm.DB) (ex error) {
		return db.Where("web_hook_id = ?", webhookId).Delete(Header{}).Error
	})
	return err
}

func (r RepositoryHeaderImpl) Save(key *tenants.TenantKey, headers Header) (Header, error) {
	e := tenants.GetConnection(key, func(db *gorm.DB) (ex error) {
		return db.Save(&headers).Error
	})
	return headers, e
}
