package accept_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/domain"
	"services.kitchen_contract/kitchen_reply"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
)

func NewAcceptTicketService(r port.Repo, sagaCmdHandler saga.CommandHandler) *AcceptTicketService {
	return &AcceptTicketService{r, sagaCmdHandler}
}

type AcceptTicketService struct {
	repo           port.Repo
	sagaCmdHandler saga.CommandHandler
}

func (s *AcceptTicketService) AcceptTicketWithOrderID(orderID int64) error {
	ticket := s.repo.FindTicketByOrderID(orderID)
	if ticket == nil {
		return apperror.New("cannot find ticket with order id " + strconv.FormatInt(orderID, 10)).
			WithCode(apperror.NotFound)
	}

	if ticket.Status == domain.TicketStatusAccepted {
		return nil
	}

	if ticket.Status != domain.TicketStatusPending {
		return apperror.New("cannot change ticket status to ACCEPTED, current status is " + ticket.Status).
			WithCode(apperror.NotAcceptable)
	}

	ticket.Status = domain.TicketStatusAccepted

	err := s.repo.UpdateTicket(ticket)
	if err != nil {
		return apperror.Wrap(err, "update ticket")
	}

	ticketCreatedReply := kitchen_reply.NewTicketCreatedReply()
	err = s.sagaCmdHandler.ReplySuccess(ticket.CommandID, ticketCreatedReply)
	if err != nil {
		return apperror.Wrap(err, "reply to command")
	}

	return nil
}
