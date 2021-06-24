package payment_command

import (
	"encoding/json"
	"services.payment_contract"
	"services.shared/saga/msg"
)

const AuthorizePayment = "AuthorizePayment"

func NewAuthorizePaymentCommand(orderID, total int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: payment_contract.PaymentServiceCommandChannel,
		Type:        AuthorizePayment,
	}
	payload := (&AuthorizePaymentPayload{
		OrderID: orderID,
		Total:   total,
	}).toJSON()
	return msg.NewCommand(meta, payload)
}

type AuthorizePaymentPayload struct {
	OrderID int64 `json:"order_id"`
	Total   int64 `json:"total"`
}

func (p *AuthorizePaymentPayload) toJSON() string {
	j, _ := json.Marshal(p)
	return string(j)
}
