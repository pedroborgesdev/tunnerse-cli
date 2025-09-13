package routes

import (
	"net/http"
	"tunnerse-server/controllers"
	"tunnerse-server/database"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := database.InitDB()
	tunnelController := controllers.NewTunnelController(db)

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	tunnel := router.Group("/")

	tunnel.POST("/open", tunnelController.Open)
	// tunnel.POST("/close", tunnelController.Get)
	// tunnel.DELETE("/delete", tunnelController.Response)
}
