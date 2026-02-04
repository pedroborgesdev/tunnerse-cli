package main

import (
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/config"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/database"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/logger"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// debug.LoadDebugConfig()

	logger.Log("INFO", "Application has been started", []logger.LogDetail{})

	config.LoadAppConfig()

	// Mostra onde os dados estão sendo salvos
	logger.Log("INFO", "Data directory", []logger.LogDetail{
		{Key: "path", Value: config.GetUserDataDir()},
		{Key: "logs", Value: config.GetLogsDir()},
		{Key: "database", Value: config.GetDatabasePath()},
	})

	database.InitDB()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":" + config.AppConfig.HTTPPort)

	select {}
}
