package appservice

import (
	"services.kitchen/internal/appservice/create_ticket"
	"services.kitchen/internal/port"
	"services.shared/saga"
)

func HandleCommands(sm saga.Manager, repo port.Repo) {
	sm.Handle(create_ticket.CreateTicketCommand, create_ticket.CreateTicketCommandHandler(repo))
}
