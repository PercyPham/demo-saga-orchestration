package approve_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/domain"
	"services.kitchen_contract/kitchen_reply"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
)

func NewApproveTicketService(repo port.Repo, sagaManager saga.Manager) *ApproveTicketService {
	return &ApproveTicketService{repo, sagaManager}
}

type ApproveTicketService struct {
	repo        port.Repo
	sagaManager saga.Manager
}

func (s *ApproveTicketService) ApproveTicket(orderID int64) error {
	ticket := s.repo.FindTicketByOrderID(orderID)
	if ticket == nil {
		return apperror.New("cannot find ticket with order id " + strconv.FormatInt(orderID, 10))
	}
	ticket.Status = domain.TicketStatusApproved
	if err := s.repo.UpdateTicket(ticket); err != nil {
		return apperror.Wrap(err, "update ticket")
	}
	err := s.sagaManager.ReplySuccess(ticket.CommandID, kitchen_reply.NewTicketApprovedReply())
	if err != nil {
		return apperror.Wrap(err, "reply TicketApproved")
	}
	return nil
}
