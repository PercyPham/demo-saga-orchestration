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

	// FindTicketByOrderID finds and returns ticket with orderID in database, return nil if not found
	FindTicketByOrderID(orderID int64) *domain.Ticket
	// FindTickets finds and returns tickets
	// 	TODO: pagination by param
	FindTickets() ([]*domain.Ticket, error)

	UpdateTicket(*domain.Ticket) error
}
