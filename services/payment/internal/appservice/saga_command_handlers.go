package appservice

import (
	"services.payment/internal/appservice/authorize_payment"
	"services.payment/internal/appservice/port"
	"services.payment_contract/payment_command"
	"services.shared/saga"
)

func HandleCommands(sm saga.Manager, repo port.Repo) {
	sm.Handle(payment_command.AuthorizePayment, authorize_payment.AuthorizePaymentCommandHandler(repo, sm))
}
