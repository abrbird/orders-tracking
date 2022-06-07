package sql_repository

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
	"strings"
)

type SQLOrderHistoryRepository struct {
	store *SQLRepository
}

func (S SQLOrderHistoryRepository) CreateOrUpdate(ctx context.Context, record *models.OrderHistoryRecord) error {
	const query = `
		INSERT INTO order_history_records (
			order_id,
			status,
			confirmation
	    ) VALUES (
			$1, $2, $3
	  	) 
	  	ON CONFLICT ON CONSTRAINT order_history_records_order_id_status_key
		DO UPDATE SET (confirmation) = ($3)
		WHERE order_history_records.order_id = $1 AND order_history_records.status=$2
		RETURNING id, order_id, status, confirmation
	`

	orderHistoryRecord := &models.OrderHistoryRecord{}
	if err := S.store.dbConnectionPool.QueryRow(
		ctx,
		query,
		record.OrderId,
		record.Status,
		record.Confirmation,
	).Scan(
		&orderHistoryRecord.Id,
		&orderHistoryRecord.OrderId,
		&orderHistoryRecord.Status,
		&orderHistoryRecord.Confirmation,
	); err != nil {
		return models.UnknownError
	}
	return nil
}

func (S SQLOrderHistoryRepository) Retrieve(ctx context.Context, recordId int64) models.OrderHistoryRecordRetrieve {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderHistoryRepository) RetrieveByStatus(ctx context.Context, orderId int64, status string) models.OrderHistoryRecordRetrieve {
	var placeholders []string
	var values []interface{}

	for _, confirmation := range []string{models.InProgress, models.Confirmed} {
		placeholders = append(
			placeholders,
			fmt.Sprintf("$%d", len(placeholders)+1),
		)
		values = append(values, confirmation)
	}
	values = append(values, orderId)
	values = append(values, status)

	query := fmt.Sprintf(`
		SELECT 
    		id,
			order_id,
			status,
			confirmation
		FROM order_history_records
		WHERE order_id = $%d 
			AND status = $%d
			AND confirmation IN (%s)
		`,
		len(placeholders)+1,
		len(placeholders)+2,
		strings.Join(placeholders, ","),
	)

	orderHistoryRecord := &models.OrderHistoryRecord{}
	if err := S.store.dbConnectionPool.QueryRow(
		ctx,
		query,
		orderId,
		status,
	).Scan(
		&orderHistoryRecord.Id,
		&orderHistoryRecord.OrderId,
		&orderHistoryRecord.Status,
		&orderHistoryRecord.Confirmation,
	); err != nil {
		return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: nil, Error: models.NotFoundError}
	}
	return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: orderHistoryRecord, Error: nil}
}

func (S SQLOrderHistoryRepository) RetrieveLast(ctx context.Context, orderId int64) models.OrderHistoryRecordRetrieve {
	var placeholders []string
	var values []interface{}

	for _, confirmation := range []string{models.InProgress, models.Confirmed} {
		placeholders = append(
			placeholders,
			fmt.Sprintf("$%d", len(placeholders)+1),
		)
		values = append(values, confirmation)
	}
	values = append(values, orderId)

	query := fmt.Sprintf(`
		SELECT 
    		id,
			order_id,
			status,
			confirmation
		FROM order_history_records
		WHERE order_id = $%d 
			AND confirmation IN (%s)
		ORDER BY updated_at DESC
		LIMIT 1`,
		len(placeholders)+1,
		strings.Join(placeholders, ","),
	)

	orderHistoryRecord := &models.OrderHistoryRecord{}
	if err := S.store.dbConnectionPool.QueryRow(
		ctx,
		query,
		orderId,
	).Scan(
		&orderHistoryRecord.Id,
		&orderHistoryRecord.OrderId,
		&orderHistoryRecord.Status,
		&orderHistoryRecord.Confirmation,
	); err != nil {
		return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: nil, Error: models.NotFoundError}
	}
	return models.OrderHistoryRecordRetrieve{OrderHistoryRecord: orderHistoryRecord, Error: nil}
}

func (S SQLOrderHistoryRepository) RetrieveHistory(ctx context.Context, orderId int64) models.OrderHistoryRetrieve {
	//TODO implement me
	panic("implement me")
}
