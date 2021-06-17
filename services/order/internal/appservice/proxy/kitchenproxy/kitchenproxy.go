package kitchenproxy

import (
	"encoding/json"
	"strconv"

	"github.com/percypham/saga-go/msg"
	"services.order/internal/domain"
	"services.shared/apperror"
)

const (
	KitchenServiceCommandChannel = "KitchenService.SagaCommandChannel"

	CommandCreateTicket  = "CreateTicket"
	CommandRejectTicket  = "RejectTicket"
	CommandApproveTicket = "ApproveTicket"
)

func GenCreateTicketCommand(order *domain.Order) (msg.Command, error) {
	meta := msg.CommandMeta{
		Destination: KitchenServiceCommandChannel,
		Type:        CommandCreateTicket,
	}
	payload, err := (&createTicketPayload{
		OrderID:   order.ID,
		Vendor:    order.Vendor,
		LineItems: order.LineItems,
	}).toJSON()
	if err != nil {
		return nil, apperror.WithLog(err, "convert payload to json")
	}
	return msg.NewCommand(meta, payload), nil
}

type createTicketPayload struct {
	OrderID   int64                   `json:"order_id"`
	Vendor    string                  `json:"vendor"`
	LineItems []*domain.OrderLineItem `json:"line_items"`
}

func (p *createTicketPayload) toJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", apperror.WithLog(err, "marshal payload")
	}
	return string(j), nil
}

func GenRejectTicketCommand(orderID int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: KitchenServiceCommandChannel,
		Type:        CommandRejectTicket,
	}
	return msg.NewCommand(meta, strconv.FormatInt(orderID, 10))
}

func GenApproveTicketCommand(orderID int64) msg.Command {
	meta := msg.CommandMeta{
		Destination: KitchenServiceCommandChannel,
		Type:        CommandApproveTicket,
	}
	return msg.NewCommand(meta, strconv.FormatInt(orderID, 10))
}
