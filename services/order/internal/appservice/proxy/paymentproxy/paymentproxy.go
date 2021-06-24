package paymentproxy

import (
	"services.payment_contract/payment_command"

	"services.order/internal/domain"
	"services.shared/saga/msg"
)

func GenAuthorizePaymentCommand(order *domain.Order) msg.Command {
	return payment_command.NewAuthorizePaymentCommand(order.ID, order.Total)
}
