package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/percypham/saga-go"
	"services.order/internal/adapter/http/rest/response"
	"services.order/internal/appservice/port"
	"services.order/internal/common/config"
	"services.shared/logger"
)

func NewOrderRestApiServer(log logger.Logger, repo port.Repo, sagaManager saga.Manager) *OrderRestApiServer {
	responder := response.New(log)
	return &OrderRestApiServer{responder, repo, sagaManager}
}

type OrderRestApiServer struct {
	response    response.Responder
	repo        port.Repo
	sagaManager saga.Manager
}

func (s *OrderRestApiServer) Run() error {
	r := gin.Default()
	if config.App().ENV != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	api := r.Group("/api")
	s.addRouteHandlers(api)
	return r.Run(":" + strconv.Itoa(config.App().PORT))
}

func (s *OrderRestApiServer) addRouteHandlers(api *gin.RouterGroup) {
	api.GET("/health", s.checkHealth)

	api.POST("/orders", s.createOrder)
	api.GET("/orders/:id", s.findOrderByID)
}
