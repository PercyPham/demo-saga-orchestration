package reject_order

import (
	"services.order/internal/appservice/port"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
)

func RejectOrderCommandHandler(repo port.Repo, ch saga.CommandHandler) func(saga.HandlerContext) error {
	return func(c saga.HandlerContext) error {
		rejectOrderService := NewRejectOrderService(repo)
		orderID, err := strconv.ParseInt(c.Command.Payload(), 10, 64)
		if err != nil {
			return apperror.Wrap(err, "get orderID from payload")
		}

		reply, err := rejectOrderService.RejectOrder(orderID)
		if err != nil {
			return apperror.Wrap(err, "reject order")
		}

		if err := c.ReplySuccess(reply); err != nil {
			return apperror.Wrap(err, "send reply")
		}

		return nil
	}
}
