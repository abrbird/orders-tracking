package implemented_service

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
	rpstr "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
)

type OrderHistoryService struct{}

func (o OrderHistoryService) Update(ctx context.Context, repository rpstr.OrderHistoryRepository, record *models.OrderHistoryRecord) error {
	return repository.CreateOrUpdate(ctx, record)
}

func (o OrderHistoryService) CreateOrUpdate(ctx context.Context, repository rpstr.OrderHistoryRepository, record *models.OrderHistoryRecord) error {
	return repository.CreateOrUpdate(ctx, record)
}

func (o OrderHistoryService) Retrieve(ctx context.Context, repository rpstr.OrderHistoryRepository, recordId int64) models.OrderHistoryRecordRetrieve {
	//TODO implement me
	panic("implement me")
}

func (o OrderHistoryService) RetrieveByStatus(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64, status string) models.OrderHistoryRecordRetrieve {
	return repository.RetrieveByStatus(ctx, orderId, status)
}

func (o OrderHistoryService) RetrieveLast(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64) models.OrderHistoryRecordRetrieve {
	return repository.RetrieveLast(ctx, orderId)
}

func (o OrderHistoryService) RetrieveHistory(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64) models.OrderHistoryRetrieve {
	//TODO implement me
	panic("implement me")
}
