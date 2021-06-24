package orderproxy

import (
	"services.order_contract/order_command"
	"services.shared/saga/msg"
)


func GenRejectOrderCommand(orderID int64) msg.Command {
	return order_command.NewRejectOrderCommand(orderID)
}

func GenApproveOrderCommand(orderID int64) msg.Command {
	return order_command.NewApproveOrderCommand(orderID)
}
