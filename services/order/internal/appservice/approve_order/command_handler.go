package approve_order

import (
	"services.order/internal/appservice/port"
	"services.shared/apperror"
	"services.shared/saga/msg"
	"strconv"
)

func ApproveOrderCommandHandler(repo port.Repo) func(command msg.Command) error {
	return func(command msg.Command) error {
		approveOrderService := NewApproveOrderService(repo)
		orderID, err := strconv.ParseInt(command.Payload(), 10, 64)
		if err != nil {
			return apperror.WithLog(err, "get orderID from payload")
		}

		err = approveOrderService.ApproveOrder(orderID)
		if err != nil {
			return apperror.WithLog(err, "approve order")
		}

		return nil
	}
}
