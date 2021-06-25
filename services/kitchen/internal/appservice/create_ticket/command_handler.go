package create_ticket

import (
	"encoding/json"
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/domain"
	"services.kitchen_contract/kitchen_command"

	"services.shared/apperror"
	"services.shared/saga/msg"
)

func CreateTicketCommandHandler(repo port.Repo) func(command msg.Command) error {
	return func(command msg.Command) error {
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
		CommandID: command.ID(),
		OrderID:   payload.OrderID,
		Vendor:    payload.Vendor,
		LineItems: lineItems,
	}, nil
}
