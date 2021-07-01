package kitchenproxy

import (
	"services.kitchen_contract/kitchen_command"
	"services.order/internal/domain"
	"services.shared/apperror"
	"services.shared/saga/msg"
)

func GenCreateTicketCommand(order *domain.Order) (msg.Command, error) {
	lineItems := make([]kitchen_command.OrderLineItem, len(order.LineItems))
	for i, item := range order.LineItems {
		lineItems[i] = kitchen_command.OrderLineItem{
			ID:   item.ID,
			Qty:  item.Quantity,
			Note: item.Note,
		}
	}

	payload := kitchen_command.CreateTicketPayload{
		OrderID:   order.ID,
		Vendor:    order.Vendor,
		LineItems: lineItems,
	}

	command, err := kitchen_command.NewCreateTicketCommand(payload)
	if err != nil {
		return nil, apperror.Wrap(err, "create CreateTicketCommand")
	}

	return command, nil
}

func GenRejectTicketCommand(orderID int64) msg.Command {
	return kitchen_command.NewRejectTicketCommand(orderID)
}

func GenApproveTicketCommand(orderID int64) msg.Command {
	return kitchen_command.NewApproveTicketCommand(orderID)
}
