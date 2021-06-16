package rest

import (
	"github.com/gin-gonic/gin"
	"services.kitchen/internal/appservice/health_check"
)

func (s *KitchenRestApiServer) checkHealth(c *gin.Context) {
	healthCheckService := health_check.NewHealthCheckService()
	if err := healthCheckService.Check(); err != nil {
		s.response.Error(c, err)
		return
	}
	s.response.Success(c, true)
}
