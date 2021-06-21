package accept_ticket

import (
	"services.kitchen/internal/domain"
	"services.kitchen/internal/port"
	"services.kitchen_contract/kitchen_reply"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
)

func NewAcceptTicketService(r port.Repo, sagaManager saga.Manager) *AcceptTicketService {
	return &AcceptTicketService{r, sagaManager}
}

type AcceptTicketService struct {
	repo        port.Repo
	sagaManager saga.Manager
}

func (s *AcceptTicketService) AcceptTicketWithOrderID(orderID int64) error {
	ticket := s.repo.FindTicketByOrderID(orderID)
	if ticket == nil {
		return apperror.New(apperror.NotFound, "cannot find ticket with order id "+strconv.FormatInt(orderID, 10))
	}

	if ticket.Status == domain.TicketStatusAccepted {
		return nil
	}

	if ticket.Status != domain.TicketStatusPending {
		return apperror.New(apperror.NotAcceptable, "ticket's status is not in PENDING, current is "+ticket.Status)
	}

	ticket.Status = domain.TicketStatusAccepted

	err := s.repo.UpdateTicket(ticket)
	if err != nil {
		return apperror.WithLog(err, "update ticket")
	}

	ticketCreatedReply := kitchen_reply.NewTicketCreatedReply()
	err = s.sagaManager.ReplySuccess(ticket.CommandID, ticketCreatedReply)
	if err != nil {
		return apperror.WithLog(err, "reply to command")
	}

	return nil
}
