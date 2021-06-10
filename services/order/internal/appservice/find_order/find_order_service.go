package find_order

import (
	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.shared/apperror"
	"strconv"
)

func NewFindOrderService(r port.OrderRepo) *FindOrderService {
	return &FindOrderService{r}
}

type FindOrderService struct {
	orderRepo port.OrderRepo
}

func (s *FindOrderService) FindByID(id int64) (*domain.Order, error) {
	order := s.orderRepo.FindOrderByID(id)
	if order == nil {
		return nil, apperror.New(apperror.NotFound, "order with id "+strconv.FormatInt(id, 10)+" not found")
	}

	return order, nil
}
