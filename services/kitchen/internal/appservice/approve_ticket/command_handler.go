package approve_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.shared/apperror"
	"services.shared/saga"
	"strconv"
)

func ApproveTicketCommandHandler(repo port.Repo) func(saga.HandlerContext) error {
	return func(c saga.HandlerContext) error {
		approveTickerService := NewApproveTicketService(repo)
		orderID, err := strconv.ParseInt(c.Command.Payload(), 10, 64)
		if err != nil {
			return apperror.Wrap(err, "get orderID from payload")
		}

		reply, err := approveTickerService.ApproveTicket(orderID)
		if err != nil {
			return apperror.Wrap(err, "approve ticket")
		}

		if err := c.ReplySuccess(reply); err != nil {
			return apperror.Wrap(err, "send reply")
		}

		return nil
	}
}
