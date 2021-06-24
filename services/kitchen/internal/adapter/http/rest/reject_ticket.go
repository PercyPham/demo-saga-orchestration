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
		s.response.Error(c, apperror.New(apperror.BadRequest, "ticket's order id must be an integer"))
		return
	}
	rejectTicketService := reject_ticket.NewRejectTicketService(s.repo, s.sagaManager)
	if err := rejectTicketService.RejectTicketWithOrderID(orderID); err != nil {
		s.response.Error(c, apperror.WithLog(err, "reject ticket"))
		return
	}
	s.response.Success(c, true)
}
