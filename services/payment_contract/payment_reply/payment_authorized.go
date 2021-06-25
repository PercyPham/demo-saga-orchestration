package payment_reply

import "services.shared/saga/msg"

const PaymentAuthorized = "PaymentAuthorized"

func NewPaymentAuthorizedReply() msg.Reply {
	meta := msg.ReplyMeta{
		Type: PaymentAuthorized,
	}
	return msg.NewReply(meta, "")
}
