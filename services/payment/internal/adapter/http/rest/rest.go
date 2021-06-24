package rest


import (
	"strconv"

	"services.shared/rest_response"

	"github.com/gin-gonic/gin"
	"services.payment/internal/appservice/port"
	"services.payment/internal/common/config"
	"services.shared/logger"
	"services.shared/saga"
)

func NewPaymentRestApiServer(log logger.Logger, repo port.Repo, sagaManager saga.Manager) *PaymentRestApiServer {
	responder := rest_response.New(log)
	responder.SetLogTrace(config.App().ENV == "development")
	return &PaymentRestApiServer{responder, repo, sagaManager}
}

type PaymentRestApiServer struct {
	response    rest_response.Responder
	repo        port.Repo
	sagaManager saga.Manager
}

func (s *PaymentRestApiServer) Run() error {
	r := gin.Default()
	if config.App().ENV != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	api := r.Group("/payment-service/api")
	s.addRouteHandlers(api)
	return r.Run(":" + strconv.Itoa(config.App().PORT))
}

func (s *PaymentRestApiServer) addRouteHandlers(api *gin.RouterGroup) {
	api.GET("/health", s.checkHealth)
}

