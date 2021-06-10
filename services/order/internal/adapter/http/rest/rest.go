package rest

import (
	"github.com/gin-gonic/gin"
	"services.order/internal/adapter/http/rest/response"
	"services.order/internal/appservice/port"
	"services.order/internal/common/config"
	"services.shared/logger"
	"strconv"
)

func NewOrderRestApiServer(log logger.Logger, repo port.Repo) *OrderRestApiServer {
	responder := response.New(log)
	return &OrderRestApiServer{ responder, repo}
}

type OrderRestApiServer struct {
	response response.Responder
	repo port.Repo
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
}
