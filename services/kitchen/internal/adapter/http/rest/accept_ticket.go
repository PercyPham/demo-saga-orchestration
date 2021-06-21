package rest

import (
	"github.com/gin-gonic/gin"
	"services.kitchen/internal/appservice/accept_ticket"
	"services.shared/apperror"
	"strconv"
)

func (s *KitchenRestApiServer) acceptTicket(c *gin.Context) {
	orderIDStr := c.Param("orderID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		s.response.Error(c, apperror.New(apperror.BadRequest, "ticket's order id must be an integer"))
		return
	}
	acceptTicketService := accept_ticket.NewAcceptTicketService(s.repo, s.sagaManager)
	if err := acceptTicketService.AcceptTicketWithOrderID(orderID); err != nil {
		s.response.Error(c, apperror.WithLog(err, "accept ticket"))
		return
	}
	s.response.Success(c, true)
}
