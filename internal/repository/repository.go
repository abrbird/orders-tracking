package repository

type Repository interface {
	OrderHistory() OrderHistoryRepository
}
