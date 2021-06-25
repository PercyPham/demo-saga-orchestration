package order_reply

import "services.shared/saga/msg"

const OrderApproved = "OrderApproved"

func NewOrderApprovedReply() msg.Reply {
	meta := msg.ReplyMeta{
		Type: OrderApproved,
	}
	return msg.NewReply(meta, "")
}

