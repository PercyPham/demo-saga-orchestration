package kitchen_command

import (
	"encoding/json"
	"services.kitchen_contract"
	"services.shared/apperror"
	"services.shared/saga/msg"
)

const CreateTicket = "CreateTicket"

type CreateTicketPayload struct {
	OrderID   int64           `json:"order_id"`
	Vendor    string          `json:"vendor"`
	LineItems []OrderLineItem `json:"line_items"`
}

type OrderLineItem struct {
	ID   string `json:"id"`
	Qty  int    `json:"qty"`
	Note string `json:"note"`
}

func (p *CreateTicketPayload) ToJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", apperror.WithLog(err, "marshal payload")
	}
	return string(j), nil
}

func NewCreateTicketCommand(payload CreateTicketPayload) (msg.Command, error) {
	meta := msg.CommandMeta{
		Destination: kitchen_contract.KitchenServiceCommandChannel,
		Type:        CreateTicket,
	}
	jsonPayload, err := payload.ToJSON()
	if err != nil {
		return nil, apperror.WithLog(err, "convert payload to json")
	}
	return msg.NewCommand(meta, jsonPayload), nil
}
