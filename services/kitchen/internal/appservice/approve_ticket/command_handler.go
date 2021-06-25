package approve_ticket

import (
	"services.kitchen/internal/appservice/port"
	"services.shared/apperror"
	"services.shared/saga"
	"services.shared/saga/msg"
	"strconv"
)

func ApproveTicketCommandHandler(repo port.Repo, sm saga.Manager) func(command msg.Command) error {
	return func(command msg.Command) error {
		approveTickerService := NewApproveTicketService(repo, sm)
		orderID, err := strconv.ParseInt(command.Payload(), 10, 64)
		if err != nil {
			return apperror.WithLog(err, "get orderID from payload")
		}

		err = approveTickerService.ApproveTicket(orderID)
		if err != nil {
			return apperror.WithLog(err, "approve ticket")
		}

		return nil
	}
}
