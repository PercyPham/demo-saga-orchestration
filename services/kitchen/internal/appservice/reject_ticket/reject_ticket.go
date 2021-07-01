package reject_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/domain"
	"services.kitchen_contract/kitchen_reply"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
)

func NewRejectTicketService(r port.Repo, sagaCmdHandler saga.CommandHandler) *RejectTicketService {
	return &RejectTicketService{r, sagaCmdHandler}
}

type RejectTicketService struct {
	repo           port.Repo
	sagaCmdHandler saga.CommandHandler
}

func (s *RejectTicketService) RejectTicketWithOrderID(orderID int64) error {
	ticket := s.repo.FindTicketByOrderID(orderID)
	if ticket == nil {
		return apperror.New("cannot find ticket with order id " + strconv.FormatInt(orderID, 10)).WithCode(apperror.NotFound)
	}

	if ticket.Status == domain.TicketStatusRejected {
		return nil
	}

	if !(ticket.Status == domain.TicketStatusPending || ticket.Status == domain.TicketStatusAccepted) {
		return apperror.New("cannot reject ticket, current status is " + ticket.Status).
			WithCode(apperror.NotAcceptable)
	}

	ticket.Status = domain.TicketStatusRejected
	if err := s.repo.UpdateTicket(ticket); err != nil {
		return apperror.Wrap(err, "update ticket")
	}

	ticketRejectedReply := kitchen_reply.NewTicketRejectedReply()
	err := s.sagaCmdHandler.ReplyFailure(ticket.CommandID, ticketRejectedReply)
	if err != nil {
		return apperror.Wrap(err, "reply to command")
	}

	return nil
}
