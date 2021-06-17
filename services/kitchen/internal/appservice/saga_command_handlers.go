package appservice

import (
	"github.com/percypham/saga-go"
	"services.kitchen/internal/appservice/create_ticket"
	"services.kitchen/internal/port"
)

func HandleCommands(sm saga.Manager, repo port.Repo) {
	sm.Handle(create_ticket.CreateTicketCommand, create_ticket.CreateTicketCommandHandler(repo))
}