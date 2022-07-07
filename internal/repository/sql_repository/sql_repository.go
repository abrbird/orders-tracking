package sql_repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
)

type SQLRepository struct {
	dbConnectionPool       *pgxpool.Pool
	orderHistoryRepository *SQLOrderHistoryRepository
}

func New(dbConnPool *pgxpool.Pool) *SQLRepository {
	repo := &SQLRepository{
		dbConnectionPool: dbConnPool,
	}
	repo.orderHistoryRepository = &SQLOrderHistoryRepository{store: repo}
	return repo
}

func (s *SQLRepository) OrderHistory() repository.OrderHistoryRepository {
	return s.orderHistoryRepository
}
