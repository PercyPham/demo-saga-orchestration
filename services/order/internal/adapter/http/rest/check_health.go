package rest

import (
	"github.com/gin-gonic/gin"
	"services.order/internal/appservice/health_check"
)

func (s *OrderRestApiServer) checkHealth(c *gin.Context) {
	healthCheckService := health_check.NewHealthCheckService(s.repo)
	if err := healthCheckService.Check(); err != nil {
		s.response.Error(c, err)
		return
	}
	s.response.Success(c, true)
}
