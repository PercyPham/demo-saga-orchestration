package appservice

import (
	"services.payment/internal/appservice/authorize_payment"
	"services.payment/internal/appservice/port"
	"services.payment_contract/payment_command"
	"services.shared/saga"
)

func HandleCommands(ch saga.CommandHandler, repo port.Repo) {
	ch.Handle(payment_command.AuthorizePayment, authorize_payment.AuthorizePaymentCommandHandler(repo, ch))
}
