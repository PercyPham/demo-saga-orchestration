package find_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/domain"
	"services.shared/apperror"
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
		return nil, apperror.Newf("ticket with order id %d not found", orderID).WithCode(apperror.NotFound)
	}
	return ticket, nil
}

func (s *FindTicketService) FindAll() ([]*domain.Ticket, error) {
	return s.ticketRepo.FindTickets()
}
