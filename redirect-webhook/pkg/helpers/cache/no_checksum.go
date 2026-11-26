package cache

import (
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/codec"
	"github.com/eko/gocache/store"
	"time"
)

type CacheWithoutChecksumKeys struct {
	codec codec.CodecInterface
}

// New instanciates a new cache entry
func NewCacheWithoutChecksumKeys(store store.StoreInterface) cache.SetterCacheInterface {
	return &CacheWithoutChecksumKeys{
		codec: codec.New(store),
	}
}

// Get returns the object stored in cache if it exists
func (c *CacheWithoutChecksumKeys) Get(key interface{}) (interface{}, error) {
	return c.codec.Get(key)
}

// GetWithTTL returns the object stored in cache and its corresponding TTL if it exists
func (c *CacheWithoutChecksumKeys) GetWithTTL(key interface{}) (interface{}, time.Duration, error) {
	return c.codec.GetWithTTL(key)
}

// Set populates the cache item using the given key
func (c *CacheWithoutChecksumKeys) Set(key, object interface{}, options *store.Options) error {
	return c.codec.Set(key, object, options)
}

// Delete removes the cache item using the given key
func (c *CacheWithoutChecksumKeys) Delete(key interface{}) error {
	return c.codec.Delete(key)
}

// Invalidate invalidates cache item from given options
func (c *CacheWithoutChecksumKeys) Invalidate(options store.InvalidateOptions) error {
	return c.codec.Invalidate(options)
}

// Clear resets all cache data
func (c *CacheWithoutChecksumKeys) Clear() error {
	return c.codec.Clear()
}

// GetCodec returns the current codec
func (c *CacheWithoutChecksumKeys) GetCodec() codec.CodecInterface {
	return c.codec
}

// GetType returns the cache type
func (c *CacheWithoutChecksumKeys) GetType() string {
	return "CacheWithoutChecksumKeys"
}
