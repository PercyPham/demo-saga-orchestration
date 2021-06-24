package order_command

import (
	"services.order_contract"
	"services.shared/saga/msg"
	"strconv"
)

const ApproveOrder = "ApproveOrder"

func NewApproveOrderCommand(orderID int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: order_contract.OrderServiceCommandChannel,
		Type:        ApproveOrder,
	}
	return msg.NewCommand(meta, strconv.FormatInt(orderID, 10))
}
