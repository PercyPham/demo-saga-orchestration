package port

import (
	"services.kitchen/internal/domain"
	"services.shared/saga"
)

type Repo interface {
	saga.Repo
	TicketRepo
}

type TicketRepo interface {
	CreateTicket(ticket *domain.Ticket) error
	FindTicketByOrderID(orderID int64) *domain.Ticket
}
