package conf

import (
	internalCacheLayer "conductor.com.br/pier/hooks/pkg/helpers/cache"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"crypto/tls"
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/metrics"
	"github.com/eko/gocache/store"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func CreateCacheManager() cache.CacheInterface {
	var applicationCacheStore store.StoreInterface
	opts := &store.Options{
		Expiration: viper.GetDuration("cache.defaultExpirationSeconds") * time.Second,
	}
	if viper.GetBool("cache.redis.enabled") {
		applicationCacheStore = createRedisStore(opts)
	} else {
		applicationCacheStore = createRistrettoStore(opts)
	}

	promMetrics := metrics.NewPrometheus(viper.GetString("app.serviceName"))
	return cache.NewMetric(promMetrics, internalCacheLayer.NewCacheWithoutChecksumKeys(applicationCacheStore))
}

func ParseCacheManagerToMarshaler(manager cache.CacheInterface) internalCacheLayer.MarshalerInterface {
	return internalCacheLayer.NewMarshaler(manager)
}

func createRistrettoStore(opts *store.Options) store.StoreInterface {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return store.NewRistretto(ristrettoCache, opts)
}

func createRedisStore(opts *store.Options) store.StoreInterface {
	var tlsConfig *tls.Config
	if viper.GetBool("cache.redis.useTls") {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("cache.redis.addr"),
		DB:           viper.GetInt("cache.redis.db"),
		Password:     viper.GetString("cache.redis.password"),
		TLSConfig:    tlsConfig,
		ReadTimeout:  time.Millisecond * time.Duration(viper.GetInt("cache.redis.readTimeout")),
		WriteTimeout: time.Millisecond * time.Duration(viper.GetInt("cache.redis.writeTimeout")),
	})

	if err := client.Ping().Err(); err != nil {
		logging.Logger.Error("[REDIS] could not ping server", zap.Error(err))
	}

	return store.NewRedis(client, opts)
}
