package redis_cache

import (
	"fmt"
	"github.com/abrbird/orders-tracking/config"
	"github.com/abrbird/orders-tracking/internal/cache"
	"time"
)

const (
	OrderHistoryRecordPrefix = "order_history_record"
)

const (
	Nil = redis.Nil
)

type RedisCache struct {
	client       *redis.Client
	orderHistory *RedisOrderHistoryCache
}

func New(cfg config.Redis) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0, // use default DB
	})

	redisCache := &RedisCache{
		client: rdb,
	}

	redisCache.orderHistory = &RedisOrderHistoryCache{
		Prefix:     OrderHistoryRecordPrefix,
		redisCache: redisCache,
		expiration: time.Minute * 15,
	}

	return redisCache
}

func (r RedisCache) OrderHistory() cache.OrderHistoryCache {
	return r.orderHistory
}
