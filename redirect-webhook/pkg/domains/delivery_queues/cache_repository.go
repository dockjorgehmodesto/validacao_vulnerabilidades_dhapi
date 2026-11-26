package delivery_queues

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/cache"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/eko/gocache/store"
	"github.com/spf13/viper"
)

const DomainCacheTag = "delivery_queue"

func NewCacheRepository(childRepo Repository, cache cache.MarshalerInterface) Repository {
	if viper.GetBool("cache.enabled") && viper.GetBool("cache.entities.queueDelivery") {
		return cacheRepository{db: childRepo, cache: cache}
	} else {
		return childRepo
	}
}

type cacheRepository struct {
	db    Repository
	cache cache.MarshalerInterface
}

// adicionando metodos
func (c cacheRepository) FindById(tenantKey *tenants.TenantKey, id string) (DeliveryQueue, error) {
	return c.db.FindById(tenantKey, id)
}

func (c cacheRepository) FindAllWithPagination(tenantKey *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error) {
	return c.db.FindAllWithPagination(tenantKey, page, limit, sort, example)
}

func (c cacheRepository) Save(tenantKey *tenants.TenantKey, deliveryQueue DeliveryQueue) (DeliveryQueue, error) {
	save, err := c.db.Save(tenantKey, deliveryQueue)
	if err == nil {
		cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId, deliveryQueue.EventType))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE DELIVERY_QUEUE] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return save, err
}

func (c cacheRepository) PartialUpdate(tenantKey *tenants.TenantKey, deliveryQueueId string, updateValues map[string]interface{}) error {
	err := c.db.PartialUpdate(tenantKey, deliveryQueueId, updateValues)
	if err == nil {
		var cacheErr error
		if _, exists := updateValues["event_type"]; exists {
			cacheErr = c.cache.Delete(GetCacheKey(tenantKey.IssuerId, updateValues["event_type"].(string)))
		} else {
			cacheErr = c.cache.Invalidate(store.InvalidateOptions{Tags: []string{FormatIssuerCacheTag(tenantKey.IssuerId)}})
		}
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE DELIVERY_QUEUE] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return err
}

func (c cacheRepository) FindByEventTypeAndIssuerId(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error) {
	queue, err := c.db.FindByEventTypeAndIssuerId(tenantKey, eventType, issuerId)
	return queue, err
}

func (c cacheRepository) Delete(tenantKey *tenants.TenantKey, id string, eventType string) error {
	err := c.db.Delete(tenantKey, id, eventType)
	if err == nil {
		cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId, eventType))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE DELIVERY_QUEUE] could not invalidate cache", zap.Error(cacheErr))
		}
	}
	return err
}

func (c cacheRepository) DeleteAll(tenantKey *tenants.TenantKey, issuerId int) error {
	err := c.db.DeleteAll(tenantKey, issuerId)
	if err == nil {
		cacheErrs := make([]error, 0)
		for _, eventType := range EventTypeList {
			cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId, eventType))
			if cacheErr != nil {
				cacheErrs = append(cacheErrs, cacheErr)
			}
		}
		if len(cacheErrs) > 0 {
			logging.Logger.Warn("[CACHE DELIVERY_QUEUE] could not invalidate some caches", zap.Errors("errors", cacheErrs))
		}
	}
	return err
}

func (c cacheRepository) DeleteIfExistsAndSaveAll(tenantKey *tenants.TenantKey, deliveryQueueList []DeliveryQueue) ([]DeliveryQueue, error) {
	queues, err := c.db.DeleteIfExistsAndSaveAll(tenantKey, deliveryQueueList)
	if err == nil {
		cacheErrs := make([]error, 0)
		for _, deliveryQueue := range deliveryQueueList {
			cacheErr := c.cache.Delete(GetCacheKey(tenantKey.IssuerId, deliveryQueue.EventType))
			if cacheErr != nil {
				cacheErrs = append(cacheErrs, cacheErr)
			}
		}
		if len(cacheErrs) > 0 {
			logging.Logger.Warn("[CACHE DELIVERY_QUEUE] could not invalidate some caches", zap.Errors("errors", cacheErrs))
		}
	}
	return queues, err
}

func (c cacheRepository) FindByEventTypeAndIssuerIdAndEnabled(tenantKey *tenants.TenantKey, eventType string, issuerId int) (DeliveryQueue, error) {
	cacheKey := GetCacheKey(issuerId, eventType)

	out := DeliveryQueue{}
	cacheDeliveryQueue, err := c.cache.Get(cacheKey, &out)
	if err == nil && cacheDeliveryQueue != nil {
		if value, ok := cacheDeliveryQueue.(*DeliveryQueue); ok {
			return *value, nil
		}
	}
	//else if err != nil {
	//	ddTrace.SetTag("cacheError", "true")
	//}

	queue, err := c.db.FindByEventTypeAndIssuerIdAndEnabled(tenantKey, eventType, issuerId)

	if err != nil && queue.Id == "" {
		queue.IssuerId = issuerId
		queue.EventType = eventType
		queue.ExchangeName = viper.GetString("rabbit.defaultExchangeName")
		queue.QueueName = viper.GetString("rabbit.defaultQueueName")
		queue.QueueKey = viper.GetString("rabbit.defaultQueueKey")

		cacheErr := c.cache.Set(cacheKey, queue, DefaultOptions(issuerId))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE DELIVERY_QUEUE] could not set cache", zap.Error(cacheErr))
		}
	}

	if err == nil {
		cacheErr := c.cache.Set(cacheKey, queue, DefaultOptions(issuerId))
		if cacheErr != nil {
			logging.Logger.Warn("[CACHE DELIVERY_QUEUE] could not set cache", zap.Error(cacheErr))
		}
	}
	return queue, err
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
	return fmt.Sprintf("delivery_queue_enabled_for_issuer_%d_%s", issuerId, eventType)
}
