package kitchen_command

import (
	"services.kitchen_contract"
	"services.shared/saga/msg"
	"strconv"
)

const RejectTicket = "RejectTicket"

func NewRejectTicketCommand(orderID int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: kitchen_contract.KitchenServiceCommandChannel,
		Type:        RejectTicket,
	}
	return msg.NewCommand(meta, strconv.FormatInt(orderID, 10))
}
