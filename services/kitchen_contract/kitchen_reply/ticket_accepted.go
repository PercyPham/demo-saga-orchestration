package kitchen_reply

import "services.shared/saga/msg"

const TicketAccepted = "TicketAccepted"

func NewTicketCreatedReply() msg.Reply {
	meta := msg.ReplyMeta{
		Type: TicketAccepted,
	}
	return msg.NewReply(meta, "")
}
