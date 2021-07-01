package appservice

import (
	"services.kitchen/internal/appservice/approve_ticket"
	"services.kitchen/internal/appservice/create_ticket"
	"services.kitchen/internal/appservice/port"
	"services.kitchen_contract/kitchen_command"
	"services.shared/saga"
)

func HandleCommands(ch saga.CommandHandler, repo port.Repo) {
	ch.Handle(kitchen_command.CreateTicket, create_ticket.CreateTicketCommandHandler(repo))
	ch.Handle(kitchen_command.ApproveTicket, approve_ticket.ApproveTicketCommandHandler(repo))
}
