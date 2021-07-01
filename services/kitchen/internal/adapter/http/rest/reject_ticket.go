package rest

import (
	"github.com/gin-gonic/gin"
	"services.kitchen/internal/appservice/reject_ticket"
	"services.shared/apperror"
	"strconv"
)

func (s *KitchenRestApiServer) rejectTicket(c *gin.Context) {
	orderIDStr := c.Param("orderID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		appErr := apperror.Wrap(err, "parse orderID").
			WithCode(apperror.BadRequest).
			WithPublicMessage("ticket's order id must be an integer")
		s.response.Error(c, appErr)
		return
	}
	rejectTicketService := reject_ticket.NewRejectTicketService(s.repo, s.sagaManager)
	if err := rejectTicketService.RejectTicketWithOrderID(orderID); err != nil {
		s.response.Error(c, apperror.Wrap(err, "reject ticket"))
		return
	}
	s.response.Success(c, true)
}
