package rest

import (
	"github.com/gin-gonic/gin"
	"services.order/internal/appservice/find_order"
	"services.shared/apperror"
	"strconv"
)

func (s *OrderRestApiServer) findOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.response.Error(c, apperror.New(apperror.BadRequest, "order id must be an integer"))
		return
	}

	findOrderService := find_order.NewFindOrderService(s.repo)
	order, err := findOrderService.FindByID(id)
	if err != nil {
		s.response.Error(c, err)
		return
	}
	s.response.Success(c, order)
}
