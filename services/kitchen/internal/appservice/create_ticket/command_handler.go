package create_ticket

import (
	"encoding/json"
	"fmt"

	"services.kitchen/internal/domain"
	"services.kitchen/internal/port"
	"services.shared/apperror"
	"services.shared/saga/msg"
)

const CreateTicketCommand = "CreateTicket"

type CommandRepo interface {
}

func CreateTicketCommandHandler(repo port.Repo) func(command msg.Command) error {
	return func(command msg.Command) error {
		if command.Type() != CreateTicketCommand {
			errMsg := fmt.Sprintf("set up wrong handler for %s command, got %s handler", command.Type(), CreateTicketCommand)
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
	payload := new(createTicketPayload)
	err := json.Unmarshal([]byte(command.Payload()), payload)
	if err != nil {
		return CreateTicketInput{}, apperror.WithLog(err, "unmarshal CreateTicket command payload")
	}
	return CreateTicketInput{
		OrderID:   payload.OrderID,
		Vendor:    payload.Vendor,
		SagaID:    command.SagaID(),
		CommandID: command.ID(),
		LineItems: payload.LineItems,
	}, nil
}

type createTicketPayload struct {
	OrderID   int64             `json:"order_id"`
	Vendor    string            `json:"vendor"`
	LineItems []domain.LineItem `json:"line_items"`
}
