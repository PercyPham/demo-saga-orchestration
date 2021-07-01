package approve_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/domain"
	"services.kitchen_contract/kitchen_reply"
	"services.shared/apperror"
	"services.shared/saga/msg"
	"strconv"
)

func NewApproveTicketService(repo port.Repo) *ApproveTicketService {
	return &ApproveTicketService{repo}
}

type ApproveTicketService struct {
	repo        port.Repo
}

func (s *ApproveTicketService) ApproveTicket(orderID int64) (msg.Reply, error) {
	ticket := s.repo.FindTicketByOrderID(orderID)
	if ticket == nil {
		return nil, apperror.New("cannot find ticket with order id " + strconv.FormatInt(orderID, 10))
	}
	ticket.Status = domain.TicketStatusApproved
	if err := s.repo.UpdateTicket(ticket); err != nil {
		return nil, apperror.Wrap(err, "update ticket")
	}
	return kitchen_reply.NewTicketApprovedReply(), nil
}
