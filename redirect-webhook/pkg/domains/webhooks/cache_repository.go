package webhooks

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/cache"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"context"
	"fmt"
	"github.com/eko/gocache/store"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"time"
)

const DomainCacheTag = "webhooks"

func NewCacheRepository(childRepo Repository, cache cache.MarshalerInterface) Repository {
	if viper.GetBool("cache.enabled") && viper.GetBool("cache.entities.webHooks") {
		return cacheRepository{db: childRepo, cache: cache}
	} else {
		return childRepo
	}
}

type cacheRepository struct {
	db    Repository
	cache cache.MarshalerInterface
}

func (c cacheRepository) FindByIdWithPreload(tenantKey *tenants.TenantKey, id string) (WebHook, error) {
	return c.db.FindByIdWithPreload(tenantKey, id)
}

func (c cacheRepository) FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, example *WebHook) (helpers.Paginator, error) {
	return c.db.FindAllWithPagination(tenantKey, page, limit, sort, example)
}

func (c cacheRepository) Save(tenantKey *tenants.TenantKey, webHook WebHook) (WebHook, error) {
	save, err := c.db.Save(tenantKey, webHook)
	if err == nil {
		cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId, webHook.EventType))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE WEBHOOKS] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return save, err
}

func (c cacheRepository) Update(tenantKey *tenants.TenantKey, id string, updateMap map[string]interface{}) error {
	err := c.db.Update(tenantKey, id, updateMap)
	if err == nil {
		var cacheErr error
		if _, exists := updateMap["event_type"]; exists {
			cacheErr = c.cache.Delete(GetCacheKey(tenantKey.IssuerId, updateMap["event_type"].(string)))
		} else {
			cacheErr = c.cache.Invalidate(store.InvalidateOptions{Tags: []string{FormatIssuerCacheTag(tenantKey.IssuerId)}})
		}
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE WEBHOOKS] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return err
}

func (c cacheRepository) FindByType(tenantKey *tenants.TenantKey, eventType string) (WebHook, error) {
	return c.db.FindByType(tenantKey, eventType)
}

func (c cacheRepository) FindByTypeWithPreload(tenantKey *tenants.TenantKey, eventType string) (WebHook, error) {
	return c.db.FindByTypeWithPreload(tenantKey, eventType)
}

func (c cacheRepository) Delete(tenantKey *tenants.TenantKey, id string, eventType string) error {
	err := c.db.Delete(tenantKey, id, eventType)
	if err == nil {
		cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId, eventType))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE WEBHOOKS] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return err
}

func (c cacheRepository) FindAllByTypeAndIssuerId(ctx context.Context, tenantKey *tenants.TenantKey, eventType string, issuerId int) ([]WebHook, error) {
	var err error
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "findWebhooksOnCache", tracer.Tag("repository", "cache"))
	defer ddTrace.Finish(tracer.WithError(err))
	cacheKey := GetCacheKey(issuerId, eventType)

	resultType := make([]WebHook, 0)
	cachedWebhooks, err := c.cache.Get(cacheKey, &resultType)
	if err == nil && cachedWebhooks != nil {
		if value, ok := cachedWebhooks.(*[]WebHook); ok {
			return *value, nil
		}
	} else if err != nil {
		ddTrace.SetTag("cacheError", "true")
	}

	webhooks, err := c.db.FindAllByTypeAndIssuerId(ctx, tenantKey, eventType, issuerId)
	if err == nil && len(webhooks) > 0 {
		cacheErr := c.cache.Set(cacheKey, webhooks, DefaultOptions(issuerId))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE WEBHOOKS] could not set cache", zap.Error(cacheErr))
		}
	}
	return webhooks, err
}

func (c cacheRepository) FindByTypeAndUrl(ctx context.Context, key *tenants.TenantKey, eventType string, url string) (WebHook, error) {
	return c.db.FindByTypeAndUrl(ctx, key, eventType, url)
}

func DefaultOptions(issuerId int) *store.Options {
	return &store.Options{
		Expiration: time.Second * viper.GetDuration("cache.defaultExpirationSeconds"),
		Tags:       FormatTags(issuerId),
	}
}

func FormatTags(issuerId int) []string {
	return []string{FormatIssuerCacheTag(issuerId), DomainCacheTag}
}

func FormatIssuerCacheTag(issuerId int) string {
	return fmt.Sprintf("%s_issuer_%d", DomainCacheTag, issuerId)
}

func GetCacheKey(issuerId int, eventType string) string {
	return fmt.Sprintf("webhooks_list_for_issuer_%d_%s", issuerId, eventType)
}
