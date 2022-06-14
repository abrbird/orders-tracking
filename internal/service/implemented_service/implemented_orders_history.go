package implemented_service

import (
	"context"
	"errors"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/cache/redis_cache"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
	rpstr "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
)

type OrderHistoryService struct {
	service *Service
}

func (o OrderHistoryService) Update(ctx context.Context, repository rpstr.OrderHistoryRepository, record *models.OrderHistoryRecord) error {

	err := repository.CreateOrUpdate(ctx, record)
	if err != nil {
		return err
	}

	if err = o.service.cache.OrderHistory().Set(ctx, *record); err != nil {
		return err
	}

	return nil
}

func (o OrderHistoryService) Create(ctx context.Context, repository rpstr.OrderHistoryRepository, record *models.OrderHistoryRecord) error {
	err := repository.CreateOrUpdate(ctx, record)
	if err != nil {
		return err
	}

	if err = o.service.cache.OrderHistory().Set(ctx, *record); err != nil {
		return err
	}

	return nil
}

func (o OrderHistoryService) RetrieveByStatus(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64, status string) models.OrderHistoryRecordRetrieve {
	retrieved := o.service.cache.OrderHistory().GetByStatus(ctx, orderId, status)
	if errors.Is(retrieved.Error, redis_cache.Nil) {
		retrieved = repository.RetrieveByStatus(ctx, orderId, status)

		if retrieved.Error != nil {
			return retrieved
		}
	}

	if err := o.service.cache.OrderHistory().Set(ctx, *retrieved.OrderHistoryRecord); err != nil {
		retrieved.OrderHistoryRecord = nil
		retrieved.Error = err
		return retrieved
	}

	return retrieved
}

func (o OrderHistoryService) Retrieve(ctx context.Context, repository rpstr.OrderHistoryRepository, recordId int64) models.OrderHistoryRecordRetrieve {

	retrieved := o.service.cache.OrderHistory().Get(ctx, recordId)
	if errors.Is(retrieved.Error, redis_cache.Nil) {
		retrieved = repository.Retrieve(ctx, recordId)

		if retrieved.Error != nil {
			return retrieved
		}
	}

	if err := o.service.cache.OrderHistory().Set(ctx, *retrieved.OrderHistoryRecord); err != nil {
		retrieved.OrderHistoryRecord = nil
		retrieved.Error = err
		return retrieved
	}

	return retrieved
}

func (o OrderHistoryService) IssueOrder(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64) models.OrderHistoryRecordRetrieve {

	issuedOrderHistoryRecordRetrieved := o.RetrieveByStatus(
		ctx,
		repository,
		orderId,
		models.Issued,
	)

	if issuedOrderHistoryRecordRetrieved.Error == nil {
		if issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.InProgress {
			issuedOrderHistoryRecordRetrieved.Error = models.OrderIsOnIssuingError(nil)
			return issuedOrderHistoryRecordRetrieved
		}
		if issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.Confirmed {
			issuedOrderHistoryRecordRetrieved.Error = models.OrderAlreadyIssuedError(nil)
			return issuedOrderHistoryRecordRetrieved
		}
	}

	readyOrderHistoryRecordRetrieved := o.RetrieveByStatus(
		ctx,
		repository,
		orderId,
		models.ReadyForIssue,
	)

	if readyOrderHistoryRecordRetrieved.Error != nil {
		readyOrderHistoryRecordRetrieved.Error = models.NotFoundError(readyOrderHistoryRecordRetrieved.Error)
		return readyOrderHistoryRecordRetrieved
	}

	if readyOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation != models.Confirmed {
		readyOrderHistoryRecordRetrieved.Error = models.OrderIsNotReadyForIssueError(nil)
		return readyOrderHistoryRecordRetrieved
	}

	issuedOrderRecord := models.OrderHistoryRecord{
		OrderId:      readyOrderHistoryRecordRetrieved.OrderHistoryRecord.OrderId,
		Status:       models.Issued,
		Confirmation: models.InProgress,
	}

	err := o.Create(
		ctx,
		repository,
		&issuedOrderRecord,
	)
	if err != nil {
		retryError := models.NewRetryError(err)
		readyOrderHistoryRecordRetrieved.Error = retryError
		return readyOrderHistoryRecordRetrieved
	}

	readyOrderHistoryRecordRetrieved.OrderHistoryRecord = &issuedOrderRecord
	return readyOrderHistoryRecordRetrieved
}

func (o OrderHistoryService) UndoIssueOrder(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64) error {

	issuedOrderHistoryRecordRetrieved := o.RetrieveByStatus(
		ctx,
		repository,
		orderId,
		models.Issued,
	)

	if issuedOrderHistoryRecordRetrieved.Error != nil {
		retryError := models.NewRetryError(issuedOrderHistoryRecordRetrieved.Error)
		return retryError
	}

	if issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.InProgress {
		issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation = models.Declined

		err := o.Update(
			ctx,
			repository,
			issuedOrderHistoryRecordRetrieved.OrderHistoryRecord,
		)
		if err != nil {
			retryError := models.NewRetryError(err)
			return retryError
		}
	}

	return nil
}

func (o OrderHistoryService) ConfirmIssueOrder(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64) error {
	issuedOrderHistoryRecordRetrieved := o.RetrieveByStatus(
		ctx,
		repository,
		orderId,
		models.Issued,
	)

	if issuedOrderHistoryRecordRetrieved.Error != nil {
		return models.NewRetryError(issuedOrderHistoryRecordRetrieved.Error)
	}

	if issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.Confirmed {
		return models.OrderAlreadyIssuedError(nil)
	}

	confirmIssueRecord := models.OrderHistoryRecord{
		Id:           issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Id,
		OrderId:      issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.OrderId,
		Status:       issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Status,
		Confirmation: models.Confirmed,
	}

	err := o.Update(
		ctx,
		repository,
		&confirmIssueRecord,
	)
	if err != nil {
		return models.NewRetryError(issuedOrderHistoryRecordRetrieved.Error)
	}

	return nil
}

func (o OrderHistoryService) RetrieveLast(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64) models.OrderHistoryRecordRetrieve {
	return repository.RetrieveLast(ctx, orderId)
}

func (o OrderHistoryService) RetrieveHistory(ctx context.Context, repository rpstr.OrderHistoryRepository, orderId int64) models.OrderHistoryRetrieve {
	//TODO implement me
	panic("implement me")
}
