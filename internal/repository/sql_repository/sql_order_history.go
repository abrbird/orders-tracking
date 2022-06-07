package sql_repository

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
)

type SQLOrderHistoryRepository struct {
	store *SQLRepository
}

func (S SQLOrderHistoryRepository) Create(ctx context.Context, record *models.OrderHistoryRecord) error {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderHistoryRepository) Retrieve(ctx context.Context, recordId int64) models.OrderHistoryRecordRetrieve {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderHistoryRepository) RetrieveByStatus(ctx context.Context, orderId int64, status string) models.OrderHistoryRecordRetrieve {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderHistoryRepository) RetrieveLast(ctx context.Context, orderId int64) models.OrderHistoryRecordRetrieve {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderHistoryRepository) RetrieveHistory(ctx context.Context, orderId int64) models.OrderHistoryRetrieve {
	//TODO implement me
	panic("implement me")
}
