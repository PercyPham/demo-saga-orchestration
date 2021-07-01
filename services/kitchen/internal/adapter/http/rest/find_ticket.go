package rest

import (
	"github.com/gin-gonic/gin"
	"services.kitchen/internal/appservice/find_ticket"
	"services.shared/apperror"
	"strconv"
)

func (s *KitchenRestApiServer) findTicketByOrderID(c *gin.Context) {
	orderIDStr := c.Param("orderID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		appErr := apperror.Wrap(err, "parse orderID").
			WithCode(apperror.BadRequest).
			WithPublicMessage("ticket's order id must be an integer")
		s.response.Error(c, appErr)
		return
	}
	findTicketService := find_ticket.NewFindTicketService(s.repo)
	ticket, err := findTicketService.FindByOrderID(orderID)
	if err != nil {
		s.response.Error(c, apperror.Wrap(err, "find ticket by order id"))
		return
	}
	s.response.Success(c, ticket)
}

func (s *KitchenRestApiServer) findAllTickets(c *gin.Context) {
	findTicketService := find_ticket.NewFindTicketService(s.repo)
	tickets, err := findTicketService.FindAll()
	if err != nil {
		s.response.Error(c, apperror.Wrap(err, "find all tickets"))
		return
	}
	s.response.Success(c, tickets)
}
