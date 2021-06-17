package rest

import (
	"github.com/gin-gonic/gin"
	"services.order/internal/appservice/create_order"
	"services.shared/apperror"
)

func (s *OrderRestApiServer) createOrder(c *gin.Context) {
	body := create_order.CreateOrderInput{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		s.response.Error(c, apperror.New(apperror.UnprocessableEntity, "invalid json: "+err.Error()))
		return
	}

	createOrderService := create_order.NewCreateOrderService(s.repo, s.sagaManager)
	order, err := createOrderService.CreateOrder(body)
	if err != nil {
		s.response.Error(c, apperror.WithLog(err, "create order using rest api"))
		return
	}
	s.response.Success(c, order)
}
