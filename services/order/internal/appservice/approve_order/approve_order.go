package approve_order

import (
	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.shared/apperror"
	"strconv"
)

func NewApproveOrderService(repo port.Repo) *ApproveOrderService {
	return &ApproveOrderService{repo}
}

type ApproveOrderService struct {
	repo port.Repo
}

func (s *ApproveOrderService) ApproveOrder(orderID int64) error {
	order := s.repo.FindOrderByID(orderID)
	if order == nil {
		return apperror.New(apperror.NotAcceptable, "cannot find order with id "+strconv.FormatInt(orderID, 10))
	}
	order.Status = domain.OrderStatusApproved
	if err := s.repo.UpdateOrder(order); err != nil {
		return apperror.WithLog(err, "update order")
	}
	return nil
}
