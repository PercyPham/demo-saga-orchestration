package paymentproxy

import (
	"encoding/json"

	"services.order/internal/domain"
	"services.shared/apperror"
	"services.shared/saga/msg"
)

const (
	PaymentServiceCommandChannel = "PaymentService.CommandChannel"

	CommandAuthorizePayment = "AuthorizePayment"
)

func GenAuthorizePaymentCommand(order *domain.Order) (msg.Command, error) {
	meta := msg.CommandMeta{
		Destination: PaymentServiceCommandChannel,
		Type:        CommandAuthorizePayment,
	}
	payload, err := (&authorizePaymentPayload{
		OrderID: order.ID,
		Total:   order.Total,
	}).toJSON()
	if err != nil {
		return nil, apperror.WithLog(err, "generate payload for AuthorizePayment command")
	}
	return msg.NewCommand(meta, payload), nil
}

type authorizePaymentPayload struct {
	OrderID int64 `json:"order_id"`
	Total   int64 `json:"total"`
}

func (p *authorizePaymentPayload) toJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", apperror.WithLog(err, "marshal authorizePaymentPayload")
	}
	return string(j), nil
}
