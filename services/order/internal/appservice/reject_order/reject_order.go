package reject_order

import (
	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.order_contract/order_reply"
	"services.shared/apperror"
	"services.shared/saga/msg"
	"strconv"
)

func NewRejectOrderService(repo port.Repo) *RejectOrderService {
	return &RejectOrderService{repo}
}

type RejectOrderService struct {
	repo port.Repo
}

func (s *RejectOrderService) RejectOrder(orderID int64) (msg.Reply, error) {
	order := s.repo.FindOrderByID(orderID)
	if order == nil {
		return nil, apperror.New("cannot find order with id " + strconv.FormatInt(orderID, 10))
	}

	order.Status = domain.OrderStatusRejected
	if err := s.repo.UpdateOrder(order); err != nil {
		return nil, apperror.Wrap(err, "update order")
	}

	return order_reply.NewOrderRejectedReply(), nil
}
