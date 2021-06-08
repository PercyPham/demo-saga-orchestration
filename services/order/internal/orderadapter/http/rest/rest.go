package rest

import (
	"os"

	"github.com/gin-gonic/gin"
)

func RunOrderServer() {
	r := gin.Default()

	if getENV("APP_ENV", "development") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	api := r.Group("/order-service/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	r.Run(":" + getENV("ORDER_SERVICE_PORT", "5000"))
}

func getENV(env, defaultVal string) string {
	if os.Getenv(env) != "" {
		return os.Getenv(env)
	}
	return defaultVal
}
