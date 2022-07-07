package sql_repository

import (
	"github.com/abrbird/orders-tracking/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
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
