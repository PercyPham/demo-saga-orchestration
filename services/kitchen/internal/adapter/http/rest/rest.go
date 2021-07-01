package rest

import (
	"services.kitchen/internal/appservice/port"
	"services.kitchen/internal/common/config"
	"services.shared/logger"
	"services.shared/rest_response"
	"services.shared/saga"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewKitchenRestApiServer(log logger.Logger, repo port.Repo, sagaCmdHandler saga.CommandHandler) *KitchenRestApiServer {
	responder := rest_response.New(log)
	responder.SetLogTrace(config.App().ENV == "development")
	return &KitchenRestApiServer{responder, repo, sagaCmdHandler}
}

type KitchenRestApiServer struct {
	response       rest_response.Responder
	repo           port.Repo
	sagaCmdHandler saga.CommandHandler
}

func (s *KitchenRestApiServer) Run() error {
	r := gin.Default()
	if config.App().ENV != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	api := r.Group("/kitchen-service/api")
	s.addRouteHandlers(api)
	return r.Run(":" + strconv.Itoa(config.App().PORT))
}

func (s *KitchenRestApiServer) addRouteHandlers(api *gin.RouterGroup) {
	api.GET("/health", s.checkHealth)

	api.GET("/tickets", s.findAllTickets)
	api.GET("/tickets/:orderID", s.findTicketByOrderID)
	api.POST("/tickets/:orderID/accept", s.acceptTicket)
	api.POST("/tickets/:orderID/reject", s.rejectTicket)
}
