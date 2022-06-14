package cache

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
)

type Cache interface {
	OrderHistory() OrderHistoryCache
}

type OrderHistoryCache interface {
	Get(ctx context.Context, recordId int64) models.OrderHistoryRecordRetrieve
	Set(ctx context.Context, record models.OrderHistoryRecord) error
	GetByStatus(ctx context.Context, orderId int64, status string) models.OrderHistoryRecordRetrieve
	SetByStatus(ctx context.Context, record models.OrderHistoryRecord) error
}
