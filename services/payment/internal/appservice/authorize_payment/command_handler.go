package authorize_payment

import (
	"encoding/json"
	"services.payment/internal/appservice/port"
	"services.payment_contract/payment_command"
	"services.shared/apperror"
	"services.shared/saga"
	"services.shared/saga/msg"
)

func AuthorizePaymentCommandHandler(repo port.Repo, ch saga.CommandHandler) func( saga.HandlerContext) error {
	return func(c  saga.HandlerContext) error {
		service := NewAuthorizePaymentService(repo, ch)
		input, err := extractAuthorizePaymentInputFromCommand(c.Command)
		if err != nil {
			return apperror.Wrap(err, "extract AuthorizePaymentInput from command")
		}

		err = service.AuthorizePayment(input)
		if err != nil {
			return apperror.Wrap(err, "authorize payment")
		}

		return nil
	}
}

func extractAuthorizePaymentInputFromCommand(command msg.Command) (AuthorizePaymentInput, error) {
	payload := new(payment_command.AuthorizePaymentPayload)
	err := json.Unmarshal([]byte(command.Payload()), payload)
	if err != nil {
		return AuthorizePaymentInput{}, apperror.Wrap(err, "unmarshal AuthorizePayment payload")
	}
	return AuthorizePaymentInput{
		OrderID:   payload.OrderID,
		Total:     payload.Total,
		CommandID: command.ID(),
	}, nil
}
