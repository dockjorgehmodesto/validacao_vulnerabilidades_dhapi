package tenants

import (
	"conductor.com.br/pier/hooks/pkg/helpers/cache"
	"fmt"
	"github.com/eko/gocache/store"
	"github.com/spf13/viper"
	"time"
)

const DomainCacheTag = "tenants"

func NewCacheRepository(childRepo Repository, cache cache.MarshalerInterface) Repository {
	if viper.GetBool("cache.enabled") && viper.GetBool("cache.entities.tenants") {
		return cacheRepository{db: childRepo, cache: cache}
	} else {
		return childRepo
	}
}

type cacheRepository struct {
	db    Repository
	cache cache.MarshalerInterface
}

func (c cacheRepository) Create(t TenantContract) (Tenant, error) {
	create, err := c.db.Create(t)
	if err == nil {
		_ = c.cache.Invalidate(store.InvalidateOptions{Tags: []string{FormatIssuerCacheKey(int(t.GetIssuerId()))}})
	}
	return create, err
}

func (c cacheRepository) Delete(id string) error {
	issuerId := 0
	deletedTenant, err := c.db.FindById(id)
	if err == nil {
		issuerId = int(deletedTenant.IssuerID)
	}

	err = c.db.Delete(id)
	if err == nil && issuerId != 0 {
		_ = c.cache.Invalidate(store.InvalidateOptions{Tags: []string{FormatIssuerCacheKey(issuerId)}})
	}
	return err
}

func (c cacheRepository) Update(id string, t *Tenant) error {
	err := c.db.Update(id, t)
	if err == nil {
		_ = c.cache.Invalidate(store.InvalidateOptions{Tags: []string{FormatIssuerCacheKey(int(t.GetIssuerId()))}})
	}
	return err
}

func (c cacheRepository) FindById(id string) (Tenant, error) {
	cacheKey := fmt.Sprintf("tenants_id_%s", id)
	cachedTenant, err := c.cache.Get(cacheKey, new(Tenant))
	if err == nil && cachedTenant != nil {
		if value, ok := cachedTenant.(*Tenant); ok {
			return *value, nil
		}
	}

	tenant, err := c.db.FindById(id)
	if err == nil {
		_ = c.cache.Set(cacheKey, tenant, DefaultOptions(int(tenant.IssuerID)))
	}
	return tenant, err
}

func (c cacheRepository) FindByIssuerIdAndType(issuerId int, accessType string) (TenantContract, error) {
	cacheKey := fmt.Sprintf("tenants_issuer_%d_type%s", issuerId, accessType)
	cachedTenant, err := c.cache.Get(cacheKey, new(Tenant))
	if err == nil && cachedTenant != nil {
		if value, ok := cachedTenant.(*Tenant); ok {
			return value, nil
		}
	}

	tenant, err := c.db.FindByIssuerIdAndType(issuerId, accessType)
	if err == nil {
		_ = c.cache.Set(cacheKey, tenant, DefaultOptions(int(tenant.GetIssuerId())))
	}
	return tenant, err
}

func (c cacheRepository) FindAll() ([]Tenant, error) {
	return c.db.FindAll()
}

func (c cacheRepository) FindAllByPagination(page int, limit int, orderBy []string) ([]Tenant, error) {
	return c.db.FindAllByPagination(page, limit, orderBy)
}

func (c cacheRepository) Ping(tenantKey *TenantKey) (SQLResult, error) {
	return c.db.Ping(tenantKey)
}

func DefaultOptions(issuerId int) *store.Options {
	return &store.Options{
		Expiration: time.Second * viper.GetDuration("cache.defaultExpirationSeconds"),
		Tags:       FormatTags(issuerId),
	}
}

func FormatTags(issuerId int) []string {
	return []string{FormatIssuerCacheKey(issuerId), DomainCacheTag}
}

func FormatIssuerCacheKey(issuerId int) string {
	return fmt.Sprintf("%s_issuer_%d", DomainCacheTag, issuerId)
}
