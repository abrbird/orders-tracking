package implemented_service

import (
	"github.com/abrbird/orders-tracking/internal/cache"
	"github.com/abrbird/orders-tracking/internal/service"
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
