package routes

import (
	"net/http"

	"github.com/pedroborgesdev/tunnerse-cli/internal/server/controllers"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/database"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := database.InitDB()
	tunnelController := controllers.NewTunnelController(db)

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	tunnel := router.Group("/")

	tunnel.POST("/new", tunnelController.New)
	tunnel.POST("/quick", tunnelController.Quick)
	tunnel.GET("/list", tunnelController.List)
	tunnel.POST("/kill", tunnelController.Kill)
	tunnel.DELETE("/delete", tunnelController.Delete)
	tunnel.POST("/info", tunnelController.Info)


}
