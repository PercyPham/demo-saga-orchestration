package orderproxy

import (
	"github.com/percypham/saga-go/msg"
	"strconv"
)

const (
	OrderServiceCommandChannel = "OrderService.SagaCommandChannel"

	CommandRejectOrder  = "RejectOrder"
	CommandApproveOrder = "ApproveOrder"
)

func GenRejectOrderCommand(orderID int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: OrderServiceCommandChannel,
		Type:        CommandRejectOrder,
	}
	return msg.NewCommand(meta, strconv.FormatInt(orderID, 10))
}

func GenApproveOrderCommand(orderID int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: OrderServiceCommandChannel,
		Type:        CommandApproveOrder,
	}
	return msg.NewCommand(meta, strconv.FormatInt(orderID, 10))
}
