package event_controller

import (
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"fmt"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Switcher interface {
	On(issuerId int, eventName string)
	Off(issuerId int, eventName string)
	set(issuerId int, eventName string, flag bool)
}

func NewSwitcher(cache cache.CacheInterface) Switcher {
	return EventSwitch{
		cache:   cache,
		enabled: viper.GetBool("cache.redis.enabled"),
	}
}

type EventSwitch struct {
	cache   cache.CacheInterface
	enabled bool
}

func (e EventSwitch) On(issuerId int, eventName string) {
	e.set(issuerId, eventName, true)
}

func (e EventSwitch) Off(issuerId int, eventName string) {
	e.set(issuerId, eventName, false)
}

func (e EventSwitch) set(issuerId int, eventName string, flag bool) {
	if e.enabled {
		cacheKey := fmt.Sprintf("enable_%d_%s", issuerId, strings.ToLower(eventName))

		err := e.cache.Set(cacheKey, strconv.FormatBool(flag), DefaultOptions())
		if err != nil {
			logging.Logger.Warn("[EVENT CONTROLLER] could not set cache", zap.Error(err))
		}
	}
}

func DefaultOptions() *store.Options {
	return &store.Options{}
}
