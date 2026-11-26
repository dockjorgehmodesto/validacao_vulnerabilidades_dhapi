package cache

import (
	"encoding/json"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
)

type MarshalerInterface interface {
	Get(key interface{}, returnObj interface{}) (interface{}, error)
	Set(key, object interface{}, options *store.Options) error
	Delete(key interface{}) error
	Invalidate(options store.InvalidateOptions) error
	Clear() error
}

type Marshaler struct {
	cache cache.CacheInterface
}

// New creates a new marshaler that marshals/unmarshals cache values
func NewMarshaler(cache cache.CacheInterface) *Marshaler {
	return &Marshaler{
		cache: cache,
	}
}

// Get obtains a value from cache and unmarshal value with given object
func (c *Marshaler) Get(key interface{}, returnObj interface{}) (interface{}, error) {
	result, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}

	switch result.(type) {
	case []byte:
		err = json.Unmarshal(result.([]byte), returnObj)

	case string:
		err = json.Unmarshal([]byte(result.(string)), returnObj)
	}

	if err != nil {
		return nil, err
	}

	return returnObj, nil
}

// Set sets a value in cache by marshaling value
func (c *Marshaler) Set(key, object interface{}, options *store.Options) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	return c.cache.Set(key, bytes, options)
}

// Delete removes a value from the cache
func (c *Marshaler) Delete(key interface{}) error {
	return c.cache.Delete(key)
}

// Invalidate invalidate cache values using given options
func (c *Marshaler) Invalidate(options store.InvalidateOptions) error {
	return c.cache.Invalidate(options)
}

// Clear reset all cache data
func (c *Marshaler) Clear() error {
	return c.cache.Clear()
}
