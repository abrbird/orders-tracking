package implemented_service

import "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/service"

type Service struct {
	orderHistoryService *OrderHistoryService
}

func New() *Service {
	return &Service{
		orderHistoryService: &OrderHistoryService{},
	}
}

func (s *Service) OrderHistory() service.OrderHistoryService {
	return s.orderHistoryService
}
