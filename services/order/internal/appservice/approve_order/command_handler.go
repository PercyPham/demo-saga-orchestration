package approve_order

import (
	"services.order/internal/appservice/port"
	"services.shared/apperror"
	"services.shared/saga"
	"services.shared/saga/msg"
	"strconv"
)

func ApproveOrderCommandHandler(repo port.Repo, sm saga.Manager) func(command msg.Command) error {
	return func(command msg.Command) error {
		approveOrderService := NewApproveOrderService(repo, sm)
		orderID, err := strconv.ParseInt(command.Payload(), 10, 64)
		if err != nil {
			return apperror.WithLog(err, "get orderID from payload")
		}

		err = approveOrderService.ApproveOrder(orderID, command.ID())
		if err != nil {
			return apperror.WithLog(err, "approve order")
		}

		return nil
	}
}
