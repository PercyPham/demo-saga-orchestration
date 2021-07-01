package approve_order

import (
	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.order_contract/order_reply"
	"services.shared/apperror"
	"services.shared/saga/msg"
	"strconv"
)

func NewApproveOrderService(repo port.Repo) *ApproveOrderService {
	return &ApproveOrderService{repo}
}

type ApproveOrderService struct {
	repo port.Repo
}

func (s *ApproveOrderService) ApproveOrder(orderID int64) (msg.Reply, error) {
	order := s.repo.FindOrderByID(orderID)
	if order == nil {
		return nil, apperror.New("cannot find order with id " + strconv.FormatInt(orderID, 10))
	}

	order.Status = domain.OrderStatusApproved
	if err := s.repo.UpdateOrder(order); err != nil {
		return nil, apperror.Wrap(err, "update order")
	}

	return order_reply.NewOrderApprovedReply(), nil
}
