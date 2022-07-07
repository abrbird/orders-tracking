package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/abrbird/orders-tracking/internal/models"
	"time"
)

type RedisOrderHistoryCache struct {
	Prefix     string
	redisCache *RedisCache
	expiration time.Duration
}

func (r RedisOrderHistoryCache) getIDKey(id int64) string {
	return fmt.Sprintf("%s_%v", r.Prefix, id)
}

func (r RedisOrderHistoryCache) getOrderStatusKey(orderId int64, status string) string {
	return fmt.Sprintf("%s_%v_%s", r.Prefix, orderId, status)
}

func (r RedisOrderHistoryCache) Get(ctx context.Context, recordId int64) models.OrderHistoryRecordRetrieve {
	//
	//	span, ctx := opentracing.StartSpanFromContext(ctx, "cache")
	//	defer span.Finish()
	//

	item, err := r.redisCache.client.Get(ctx, r.getIDKey(recordId)).Result()
	if err != nil {
		return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: nil, Error: err}
	}

	var record models.OrderHistoryRecord
	if err = json.Unmarshal([]byte(item), &record); err != nil {
		return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: nil, Error: err}
	}

	return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: &record, Error: nil}
}

func (r RedisOrderHistoryCache) Set(ctx context.Context, record models.OrderHistoryRecord) error {
	//
	//	span, ctx := opentracing.StartSpanFromContext(ctx, "cache")
	//	defer span.Finish()
	//

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return r.redisCache.client.Set(ctx, r.getIDKey(record.Id), data, r.expiration).Err()
}

func (r RedisOrderHistoryCache) GetByStatus(ctx context.Context, orderId int64, status string) models.OrderHistoryRecordRetrieve {
	//
	//	span, ctx := opentracing.StartSpanFromContext(ctx, "cache")
	//	defer span.Finish()
	//

	item, err := r.redisCache.client.Get(ctx, r.getOrderStatusKey(orderId, status)).Result()
	if err != nil {
		return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: nil, Error: err}
	}

	var record models.OrderHistoryRecord
	if err = json.Unmarshal([]byte(item), &record); err != nil {
		return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: nil, Error: err}
	}

	return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: &record, Error: nil}
}

func (r RedisOrderHistoryCache) SetByStatus(ctx context.Context, record models.OrderHistoryRecord) error {
	//
	//	span, ctx := opentracing.StartSpanFromContext(ctx, "cache")
	//	defer span.Finish()
	//

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return r.redisCache.client.Set(ctx, r.getOrderStatusKey(record.OrderId, record.Status), data, r.expiration).Err()
}
