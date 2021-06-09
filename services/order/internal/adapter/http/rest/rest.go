package rest

import (
	"os"
	"services.order/internal/appservice/health_check"
	"services.order/internal/appservice/port"

	"github.com/gin-gonic/gin"
)

func RunOrderServer(repo port.Repo) {
	r := gin.Default()

	if getENV("APP_ENV", "development") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	api := r.Group("/order-service/api")

	api.GET("/health", func(c *gin.Context) {
		healthCheckService := health_check.NewHealthCheckService(repo)
		err := healthCheckService.Check()
		if err != nil {
			c.JSON(200, gin.H{"message":"OK"})
			return
		}
		c.JSON(200, gin.H{"error":err.Error()})
	})

	r.Run(":" + getENV("ORDER_SERVICE_PORT", "5000"))
}

func getENV(env, defaultVal string) string {
	if os.Getenv(env) != "" {
		return os.Getenv(env)
	}
	return defaultVal
}
