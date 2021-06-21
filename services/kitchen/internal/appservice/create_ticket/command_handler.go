package create_ticket

import (
	"encoding/json"
	"fmt"
	"services.kitchen/internal/domain"
	"services.kitchen_contract/kitchen_command"

	"services.kitchen/internal/port"
	"services.shared/apperror"
	"services.shared/saga/msg"
)

type CommandRepo interface {
}

func CreateTicketCommandHandler(repo port.Repo) func(command msg.Command) error {
	return func(command msg.Command) error {
		if command.Type() != kitchen_command.CreateTicket {
			errMsg := fmt.Sprintf("set up wrong handler for %s command, got %s handler",
				command.Type(), kitchen_command.CreateTicket)
			return apperror.New(apperror.InternalServerError, errMsg)
		}

		createTickerService := NewCreateTicketService(repo)
		input, err := extractCreateTicketInputFromCommand(command)
		if err != nil {
			return apperror.WithLog(err, "extract CreateTicketInput from command")
		}

		err = createTickerService.CreateTicket(input)
		if err != nil {
			return apperror.WithLog(err, "create ticket")
		}

		return nil
	}
}

type CommandReplyMeta struct {
	SagaID    string
	CommandID string
}

func extractCreateTicketInputFromCommand(command msg.Command) (CreateTicketInput, error) {
	payload := new(kitchen_command.CreateTicketPayload)
	err := json.Unmarshal([]byte(command.Payload()), payload)
	if err != nil {
		return CreateTicketInput{}, apperror.WithLog(err, "unmarshal CreateTicket command payload")
	}
	lineItems := make([]domain.LineItem, len(payload.LineItems))
	for i, item := range payload.LineItems {
		lineItems[i] = domain.LineItem{
			ID:       item.ID,
			Quantity: item.Qty,
			Note:     item.Note,
		}
	}
	return CreateTicketInput{
		OrderID:   payload.OrderID,
		Vendor:    payload.Vendor,
		SagaID:    command.SagaID(),
		CommandID: command.ID(),
		LineItems: lineItems,
	}, nil
}

