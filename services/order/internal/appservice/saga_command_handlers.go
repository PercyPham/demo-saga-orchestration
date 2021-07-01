package appservice

import (
	"services.order/internal/appservice/approve_order"
	"services.order/internal/appservice/port"
	"services.order/internal/appservice/reject_order"
	"services.order_contract/order_command"
	"services.shared/saga"
)

func HandleCommands(ch saga.CommandHandler, repo port.Repo) {
	ch.Handle(order_command.ApproveOrder, approve_order.ApproveOrderCommandHandler(repo, ch))
	ch.Handle(order_command.RejectOrder, reject_order.RejectOrderCommandHandler(repo, ch))
}
