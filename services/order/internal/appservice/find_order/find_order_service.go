package find_order

import (
	"strconv"

	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.shared/apperror"
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
		return nil, apperror.New("order with id " + strconv.FormatInt(id, 10) + " not found").WithCode(apperror.NotFound)
	}

	return order, nil
}

func (s *FindOrderService) FindAll() ([]*domain.Order, error) {
	return s.orderRepo.FindOrders()
}
