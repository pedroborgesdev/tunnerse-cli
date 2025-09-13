package main

import (
	"tunnerse-server/config"
	"tunnerse-server/database"
	"tunnerse-server/debug"
	"tunnerse-server/logger"
	"tunnerse-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	debug.LoadDebugConfig()

	logger.Log("INFO", "Application has been started", []logger.LogDetail{})

	config.LoadAppConfig()

	database.InitDB()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":" + config.AppConfig.HTTPPort)

	select {}
}
