package approve_order

import (
	"services.order/internal/appservice/port"
	"services.order/internal/domain"
	"services.order_contract/order_reply"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
	"time"
)

func NewApproveOrderService(repo port.Repo, sm saga.Manager) *ApproveOrderService {
	return &ApproveOrderService{repo, sm}
}

type ApproveOrderService struct {
	repo        port.Repo
	sagaManager saga.Manager
}

func (s *ApproveOrderService) ApproveOrder(orderID int64, commandID string) error {
	order := s.repo.FindOrderByID(orderID)
	if order == nil {
		return apperror.New(apperror.NotAcceptable, "cannot find order with id "+strconv.FormatInt(orderID, 10))
	}
	order.Status = domain.OrderStatusApproved
	if err := s.repo.UpdateOrder(order); err != nil {
		return apperror.WithLog(err, "update order")
	}

	go s.replyOrderApproved(commandID)

	return nil
}

// TODO: this is hot fix, need to have immediate saga reply, re-implement later
func (s *ApproveOrderService) replyOrderApproved(commandID string) {
	time.Sleep(1000)
	err := s.sagaManager.ReplySuccess(commandID, order_reply.NewOrderApprovedReply())
	if err != nil {
		panic(apperror.WithLog(err, "reply OrderApproved"))
	}
}
