package rest

import (
	"services.kitchen/internal/adapter/http/rest/response"
	"services.kitchen/internal/common/config"
	"services.shared/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewKitchenRestApiServer(log logger.Logger) *KitchenRestApiServer {
	responder := response.New(log)
	responder.SetLogTrace(config.App().ENV == "development")
	return &KitchenRestApiServer{responder}
}

type KitchenRestApiServer struct {
	response response.Responder
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
}
