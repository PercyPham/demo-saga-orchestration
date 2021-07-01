package order_reply

import "services.shared/saga/msg"

const OrderRejected = "OrderRejected"

func NewOrderRejectedReply() msg.Reply {
	meta := msg.ReplyMeta{
		Type: OrderRejected,
	}
	return msg.NewReply(meta, "")
}

