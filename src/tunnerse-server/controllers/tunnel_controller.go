package controllers

import (
	"tunnerse-server/database"
	"tunnerse-server/logger"
	"tunnerse-server/services"
	"tunnerse-server/utils"

	"github.com/gin-gonic/gin"
)

type TunnelController struct {
	tunnelService *services.TunnelService
}

func NewTunnelController(db *database.Database) *TunnelController {
	return &TunnelController{
		tunnelService: services.NewTunnelService(db),
	}
}

func (c *TunnelController) Open(ctx *gin.Context) {
	var req utils.OpenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		return
	}

	tunnelName, isSubdomain, err := c.tunnelService.RegisterTunnel(req.Name, req.Port, req.ServerURL)
	if err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		logger.Log("ERROR", "Registration failed", []logger.LogDetail{{Key: "Error", Value: err.Error()}})
		return
	}

	utils.Success(ctx, gin.H{
		"message":   "tunnel has been registered",
		"subdomain": isSubdomain,
		"tunnel":    tunnelName,
	})
	logger.Log("INFO", "User registered successfully", []logger.LogDetail{
		{Key: "subdomain", Value: isSubdomain},
		{Key: "tunnel", Value: tunnelName},
	})
}
