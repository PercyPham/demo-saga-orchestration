package kitchen_reply

import "services.shared/saga/msg"

const TicketRejected = "TicketRejected"

func NewTicketRejectedReply() msg.Reply {
	meta := msg.ReplyMeta{
		Type: TicketRejected,
	}
	return msg.NewReply(meta, "")
}
