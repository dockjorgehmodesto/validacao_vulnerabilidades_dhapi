package encryption

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

const DomainCacheTag = "keyvaults"

func NewCacheRepository(childRepo Repository, cache cache.MarshalerInterface) Repository {
	if viper.GetBool("cache.enabled") && viper.GetBool("cache.entities.keyVaults") {
		return cacheRepository{db: childRepo, cache: cache}
	} else {
		return childRepo
	}
}

type cacheRepository struct {
	db    Repository
	cache cache.MarshalerInterface
}

func (c cacheRepository) FindById(tenantKey *tenants.TenantKey, id string) (KeyVault, error) {
	return c.db.FindById(tenantKey, id)
}

func (c cacheRepository) FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error) {
	return c.db.FindAllWithPagination(tenantKey, page, limit, sort, example)
}

func (c cacheRepository) Save(tenantKey *tenants.TenantKey, keyvault KeyVault) (KeyVault, error) {
	save, err := c.db.Save(tenantKey, keyvault)
	if err == nil {
		cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE KEYVAULTS] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return save, err
}

func (c cacheRepository) Update(tenantKey *tenants.TenantKey, keyvault KeyVault) (KeyVault, error) {
	update, err := c.db.Update(tenantKey, keyvault)
	if err == nil {
		cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE KEYVAULTS] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return update, err
}

func (c cacheRepository) FindFirst(tenantKey *tenants.TenantKey) (KeyVault, error) {
	return c.db.FindFirst(tenantKey)
}

func (c cacheRepository) FindMostRecentAndActive(tenantKey *tenants.TenantKey) (KeyVault, error) {
	return c.db.FindMostRecentAndActive(tenantKey)
}

func (c cacheRepository) FindMostRecentAndActiveByIssuerId(ctx context.Context, tenantKey *tenants.TenantKey, issuerId int) (KeyVault, error) {
	var err error
	ddTrace, ctx := tracer.StartSpanFromContext(ctx, "findKeyvaultOnCache", tracer.Tag("repository", "cache"))
	defer ddTrace.Finish(tracer.WithError(err))
	cacheKey := GetCacheKey(issuerId)

	cachedKeyVault, err := c.cache.Get(cacheKey, new(KeyVault))
	if err == nil && cachedKeyVault != nil {
		if value, ok := cachedKeyVault.(*KeyVault); ok {
			return *value, nil
		}
	} else if err != nil {
		ddTrace.SetTag("cacheError", "true")
	}

	keyVault, err := c.db.FindMostRecentAndActiveByIssuerId(ctx, tenantKey, issuerId)
	if err == nil {
		cacheErr := c.cache.Set(cacheKey, keyVault, DefaultOptions())
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE KEYVAULTS] could not set cache", zap.Error(cacheErr))
		}
	}
	return keyVault, err

}

func DefaultOptions() *store.Options {
	return &store.Options{
		Expiration: time.Second * viper.GetDuration("cache.defaultExpirationSeconds"),
		Tags:       FormatTags(),
	}
}

func FormatTags() []string {
	return []string{DomainCacheTag}
}

func GetCacheKey(issuerId int) string {
	return fmt.Sprintf("key_vaults_most_recent_and_active_%d", issuerId)
}
