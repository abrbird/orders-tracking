package implemented_service

import (
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/cache"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/service"
)

type Service struct {
	cache               cache.Cache
	orderHistoryService *OrderHistoryService
}

func New(cache cache.Cache) *Service {
	srvc := &Service{
		cache: cache,
	}
	srvc.orderHistoryService = &OrderHistoryService{srvc}

	return srvc
}

func (s *Service) OrderHistory() service.OrderHistoryService {
	return s.orderHistoryService
}
