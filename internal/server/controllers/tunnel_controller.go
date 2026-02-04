package controllers

import (
	"strings"

	"github.com/pedroborgesdev/tunnerse-cli/internal/server/database"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/logger"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/services"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/utils"

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

func (c *TunnelController) New(ctx *gin.Context) {
	var req utils.OpenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		return
	}

	tunnelName, isSubdomain, err := c.tunnelService.RegisterTunnel(req.Name, req.Port, req.ServerURL, false)
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
	logger.Log("INFO", "Tunnel registered successfully", []logger.LogDetail{
		{Key: "subdomain", Value: isSubdomain},
		{Key: "tunnel", Value: tunnelName},
	})
}

func (c *TunnelController) Quick(ctx *gin.Context) {
	var req utils.OpenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		return
	}

	tunnelName, isSubdomain, err := c.tunnelService.RegisterTunnel(req.Name, req.Port, req.ServerURL, true)
	if err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		logger.Log("ERROR", "Quick tunnel registration failed", []logger.LogDetail{{Key: "Error", Value: err.Error()}})
		return
	}

	utils.Success(ctx, gin.H{
		"message":   "quick tunnel has been registered",
		"subdomain": isSubdomain,
		"tunnel":    tunnelName,
	})
	logger.Log("INFO", "Quick tunnel registered successfully", []logger.LogDetail{
		{Key: "subdomain", Value: isSubdomain},
		{Key: "tunnel", Value: tunnelName},
	})
}

func (c *TunnelController) List(ctx *gin.Context) {
	tunnels, err := c.tunnelService.ListTunnels()
	if err != nil {
		utils.InternalError(ctx, gin.H{"error": err.Error()})
		logger.Log("ERROR", "Failed to list tunnels", []logger.LogDetail{{Key: "Error", Value: err.Error()}})
		return
	}

	utils.Success(ctx, gin.H{
		"tunnels": tunnels,
		"count":   len(tunnels),
	})
}

func (c *TunnelController) Kill(ctx *gin.Context) {
	var req utils.KillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		return
	}

	err := c.tunnelService.KillTunnel(req.TunnelID)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "tunnel not found") {
			utils.NotFound(ctx, gin.H{"error": "tunnel not found", "tunnel_id": req.TunnelID})
			logger.Log("WARN", "Tunnel not found for kill", []logger.LogDetail{{Key: "tunnel_id", Value: req.TunnelID}})
			return
		}

		utils.InternalError(ctx, gin.H{"error": err.Error(), "tunnel_id": req.TunnelID})
		logger.Log("ERROR", "Failed to kill tunnel", []logger.LogDetail{{Key: "Error", Value: err.Error()}, {Key: "tunnel_id", Value: req.TunnelID}})
		return
	}

	utils.Success(ctx, gin.H{
		"message":   "tunnel has been killed",
		"tunnel_id": req.TunnelID,
	})
	logger.Log("INFO", "Tunnel killed successfully", []logger.LogDetail{
		{Key: "tunnel_id", Value: req.TunnelID},
	})
}

func (c *TunnelController) Delete(ctx *gin.Context) {
	var req utils.DeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		return
	}

	err := c.tunnelService.DeleteTunnel(req.TunnelID)
	if err != nil {

		errMsg := err.Error()
		if strings.Contains(errMsg, "tunnel not found") {
			utils.NotFound(ctx, gin.H{"error": "tunnel not found", "tunnel_id": req.TunnelID})
			logger.Log("WARN", "Tunnel not found for deletion", []logger.LogDetail{{Key: "tunnel_id", Value: req.TunnelID}})
			return
		}

		if strings.Contains(errMsg, "still active") {
			utils.BadRequest(ctx, gin.H{"error": "tunnel is still active, please kill it first", "tunnel_id": req.TunnelID})
			logger.Log("WARN", "Attempted to delete active tunnel", []logger.LogDetail{{Key: "tunnel_id", Value: req.TunnelID}})
			return
		}

		utils.InternalError(ctx, gin.H{"error": err.Error()})
		logger.Log("ERROR", "Failed to delete tunnel", []logger.LogDetail{{Key: "Error", Value: err.Error()}})
		return
	}

	utils.Success(ctx, gin.H{
		"message":   "tunnel has been deleted",
		"tunnel_id": req.TunnelID,
	})
	logger.Log("INFO", "Tunnel deleted successfully", []logger.LogDetail{
		{Key: "tunnel_id", Value: req.TunnelID},
	})
}

func (c *TunnelController) Info(ctx *gin.Context) {
	var req utils.DeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, gin.H{"error": err.Error()})
		return
	}

	info, err := c.tunnelService.GetTunnelInfo(req.TunnelID)
	if err != nil {

		errMsg := err.Error()
		if strings.Contains(errMsg, "tunnel not found") || strings.Contains(errMsg, "info not found") {
			utils.NotFound(ctx, gin.H{"error": "tunnel not found", "tunnel_id": req.TunnelID})
			logger.Log("WARN", "Tunnel not found for info", []logger.LogDetail{{Key: "tunnel_id", Value: req.TunnelID}})
			return
		}

		utils.InternalError(ctx, gin.H{"error": err.Error()})
		logger.Log("ERROR", "Failed to get tunnel info", []logger.LogDetail{{Key: "Error", Value: err.Error()}})
		return
	}

	utils.Success(ctx, gin.H{
		"info": info,
	})
}
