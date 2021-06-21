package kitchen_command

import (
	"services.kitchen_contract"
	"services.shared/saga/msg"
	"strconv"
)

const ApproveTicket = "ApproveTicket"

func NewApproveTicketCommand(orderID int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: kitchen_contract.KitchenServiceCommandChannel,
		Type:        ApproveTicket,
	}
	return msg.NewCommand(meta, strconv.FormatInt(orderID, 10))
}
