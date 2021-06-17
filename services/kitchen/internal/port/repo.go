package port

import (
	"github.com/percypham/saga-go"
	"services.kitchen/internal/domain"
)

type Repo interface {
	saga.Repo
	TicketRepo
}

type TicketRepo interface {
	CreateTicket(ticket *domain.Ticket) error
	FindTicketByOrderID(orderID int64) *domain.Ticket
}