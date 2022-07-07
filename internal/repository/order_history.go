package repository

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
)

type OrderHistoryRepository interface {
	CreateOrUpdate(ctx context.Context, record *models.OrderHistoryRecord) error
	Retrieve(ctx context.Context, recordId int64) models.OrderHistoryRecordRetrieve
	RetrieveByStatus(ctx context.Context, orderId int64, status string) models.OrderHistoryRecordRetrieve
	RetrieveLast(ctx context.Context, orderId int64) models.OrderHistoryRecordRetrieve
	RetrieveHistory(ctx context.Context, orderId int64) models.OrderHistoryRetrieve
}
