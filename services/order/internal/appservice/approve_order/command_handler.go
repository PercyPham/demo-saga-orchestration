package approve_order

import (
	"services.order/internal/appservice/port"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
)

func ApproveOrderCommandHandler(repo port.Repo, ch saga.CommandHandler) func(saga.HandlerContext) error {
	return func(c saga.HandlerContext) error {
		approveOrderService := NewApproveOrderService(repo)
		orderID, err := strconv.ParseInt(c.Command.Payload(), 10, 64)
		if err != nil {
			return apperror.Wrap(err, "get orderID from payload")
		}

		reply, err := approveOrderService.ApproveOrder(orderID)
		if err != nil {
			return apperror.Wrap(err, "approve order")
		}

		if err := c.ReplySuccess(reply); err != nil {
			return apperror.Wrap(err, "send reply")
		}

		return nil
	}
}
