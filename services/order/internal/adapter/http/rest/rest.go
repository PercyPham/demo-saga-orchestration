package rest

import (
	"strconv"

	"services.shared/rest_response"

	"github.com/gin-gonic/gin"
	"services.order/internal/appservice/port"
	"services.order/internal/common/config"
	"services.shared/logger"
	"services.shared/saga"
)

func NewOrderRestApiServer(log logger.Logger, repo port.Repo, sagaManager saga.Manager) *OrderRestApiServer {
	responder := rest_response.New(log)
	responder.SetLogTrace(config.App().ENV == "development")
	return &OrderRestApiServer{responder, repo, sagaManager}
}

type OrderRestApiServer struct {
	response    rest_response.Responder
	repo        port.Repo
	sagaManager saga.Manager
}

func (s *OrderRestApiServer) Run() error {
	r := gin.Default()
	if config.App().ENV != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	api := r.Group("/order-service/api")
	s.addRouteHandlers(api)
	return r.Run(":" + strconv.Itoa(config.App().PORT))
}

func (s *OrderRestApiServer) addRouteHandlers(api *gin.RouterGroup) {
	api.GET("/health", s.checkHealth)

	api.POST("/orders", s.createOrder)
	api.GET("/orders", s.findAllOrders)
	api.GET("/orders/:id", s.findOrderByID)
}
