package find_ticket

import (
	"services.kitchen/internal/domain"
	"services.kitchen/internal/port"
	"services.shared/apperror"
	"strconv"
)

func NewFindTicketService(r port.TicketRepo) *FindTicketService {
	return &FindTicketService{r}
}

type FindTicketService struct {
	ticketRepo port.TicketRepo
}

func (s *FindTicketService) FindByOrderID(orderID int64) (*domain.Ticket, error) {
	ticket := s.ticketRepo.FindTicketByOrderID(orderID)
	if ticket == nil {
		return nil, apperror.New(apperror.NotFound, "ticket with order id "+strconv.FormatInt(orderID, 10)+" not found")
	}
	return ticket, nil
}

func (s *FindTicketService) FindAll() ([]*domain.Ticket, error) {
	return s.ticketRepo.FindTickets()
}
