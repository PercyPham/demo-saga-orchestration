package kitchen_reply

import "services.shared/saga/msg"

const TicketApproved = "TicketApproved"

func NewTicketApprovedReply() msg.Reply {
	meta := msg.ReplyMeta{
		Type: TicketApproved,
	}
	return msg.NewReply(meta, "")
}
