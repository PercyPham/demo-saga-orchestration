package appservice

import (
	"services.order/internal/appservice/approve_order"
	"services.order/internal/appservice/port"
	"services.order_contract/order_command"
	"services.shared/saga"
)

func HandleCommands(sm saga.Manager, repo port.Repo) {
	sm.Handle(order_command.ApproveOrder, approve_order.ApproveOrderCommandHandler(repo, sm))
}
