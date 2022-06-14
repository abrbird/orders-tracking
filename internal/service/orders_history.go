package service

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
)

type OrderHistoryService interface {
	Create(ctx context.Context, repository repository.OrderHistoryRepository, record *models.OrderHistoryRecord) error
	Update(ctx context.Context, repository repository.OrderHistoryRepository, record *models.OrderHistoryRecord) error
	Retrieve(ctx context.Context, repository repository.OrderHistoryRepository, recordId int64) models.OrderHistoryRecordRetrieve
	RetrieveByStatus(ctx context.Context, repository repository.OrderHistoryRepository, orderId int64, status string) models.OrderHistoryRecordRetrieve
	RetrieveLast(ctx context.Context, repository repository.OrderHistoryRepository, orderId int64) models.OrderHistoryRecordRetrieve
	RetrieveHistory(ctx context.Context, repository repository.OrderHistoryRepository, orderId int64) models.OrderHistoryRetrieve

	IssueOrder(ctx context.Context, repository repository.OrderHistoryRepository, orderId int64) models.OrderHistoryRecordRetrieve
	UndoIssueOrder(ctx context.Context, repository repository.OrderHistoryRepository, orderId int64) error
	ConfirmIssueOrder(ctx context.Context, repository repository.OrderHistoryRepository, orderId int64) error
}
