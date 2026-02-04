package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pedroborgesdev/tunnerse-cli/internal/server/config"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/database"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/jobs"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/models"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/repositories"
)

type TunnelService struct {
	repo *repositories.TunnelRepository
}

func NewTunnelService(db *database.Database) *TunnelService {
	repo := repositories.NewTunnelRepository(db)
	return &TunnelService{
		repo: repo,
	}
}

func (s *TunnelService) RegisterTunnel(name, port, server_url string, isQuick bool) (string, bool, error) {
	payload := map[string]string{"name": name}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", false, fmt.Errorf("encode JSON: %w", err)
	}

	registerURL := fmt.Sprintf("%s/register", server_url)
	resp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", false, fmt.Errorf("post register: %w", err)
	}
	defer resp.Body.Close()

	var result models.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", false, fmt.Errorf("decode register response. probably tunnerse server is offline: %w", err)
	}

	tunnelFullURL := result.Data.Tunnel

	serverDomain := strings.TrimPrefix(server_url, "http://")
	serverDomain = strings.TrimPrefix(serverDomain, "https://")

	// Extrai o ID do túnel da URL
	tunnelID := extractTunnelID(tunnelFullURL, serverDomain)

	protocol := "http://"
	if strings.HasPrefix(server_url, "https://") {
		protocol = "https://"
	}

	var finalTunnelURL string
	if strings.HasPrefix(tunnelFullURL, "http://") || strings.HasPrefix(tunnelFullURL, "https://") {
		finalTunnelURL = tunnelFullURL
	} else {
		if result.Data.Subdomain {
			finalTunnelURL = fmt.Sprintf("%s%s.%s", protocol, tunnelFullURL, serverDomain)
		} else {
			finalTunnelURL = fmt.Sprintf("%s%s/%s", protocol, serverDomain, tunnelFullURL)
		}
	}

	if !isQuick {
		tunnel := &models.Tunnel{
			ID:        tunnelID,
			Port:      port,
			Url:       finalTunnelURL,
			Domain:    server_url,
			Active:    true,
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		info := &models.Info{
			ID:           tunnelID,
			Requests:     0,
			Healthchecks: 0,
			Warns:        0,
			Errors:       0,
		}

		if err := s.repo.Create(tunnel, info); err != nil {
			return "", false, fmt.Errorf("tunnel not saved: %w", err)
		}
	} else {
		config.QuickTunnelURLs[tunnelID] = finalTunnelURL
	}

	loopJob := jobs.NewLoopJob(s.repo.DB, tunnelID, port, result.Data.Subdomain, server_url, finalTunnelURL, isQuick)
	if loopJob == nil {
		return "", false, fmt.Errorf("failed to create tunnel job")
	}

	config.SetActiveJob(tunnelID, loopJob)

	go func() {
		loopJob.StartTunnelLoop()
		config.RemoveActiveJob(tunnelID)
	}()

	return tunnelID, result.Data.Subdomain, nil
}

func (s *TunnelService) ListTunnels() ([]*models.Tunnel, error) {
	return s.repo.ListTunnels()
}

func (s *TunnelService) KillTunnel(tunnelID string) error {
	var tunnelURL string
	isQuickTunnel := false

	if url, exists := config.QuickTunnelURLs[tunnelID]; exists {
		tunnelURL = url
		isQuickTunnel = true
	} else {
		tunnel, err := s.repo.GetTunnel(tunnelID)
		if err != nil {
			return fmt.Errorf("tunnel not found: %w", err)
		}
		tunnelURL = tunnel.Url
	}

	if tunnelURL == "" {
		return fmt.Errorf("tunnel URL is empty")
	}

	go func() {
		closeURL := tunnelURL + "/close"

		payload := map[string]string{"name": tunnelID}
		data, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("failed to marshal close payload: %v\n", err)
			return
		}

		resp, err := http.Post(closeURL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Printf("failed to send close request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		job, exists := config.GetActiveJob(tunnelID)
		if exists {
			config.RemoveActiveJob(tunnelID)
			job.Stop()
		}

		if isQuickTunnel {
			delete(config.QuickTunnelURLs, tunnelID)
		} else {
			if err := s.repo.UpdateTunnelStatus(tunnelID, false); err != nil {
				fmt.Printf("failed to update tunnel status: %v\n", err)
			}
		}
	}()

	return nil
}

func (s *TunnelService) DeleteTunnel(tunnelID string) error {
	tunnel, err := s.repo.GetTunnel(tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel not found: %w", err)
	}

	if tunnel.Active {
		return fmt.Errorf("tunnel is still active, please kill it first")
	}

	if err := s.repo.DeleteTunnel(tunnelID); err != nil {
		return fmt.Errorf("failed to delete tunnel: %w", err)
	}

	return nil
}

func (s *TunnelService) GetTunnelInfo(tunnelID string) (map[string]interface{}, error) {
	tunnel, err := s.repo.GetTunnel(tunnelID)
	if err != nil {
		return nil, fmt.Errorf("tunnel not found: %w", err)
	}

	info, err := s.repo.GetInfo(tunnelID)
	if err != nil {
		return nil, fmt.Errorf("info not found: %w", err)
	}

	result := map[string]interface{}{
		"id":           tunnel.ID,
		"port":         tunnel.Port,
		"url":          tunnel.Url,
		"domain":       tunnel.Domain,
		"active":       tunnel.Active,
		"created_at":   tunnel.CreatedAt,
		"requests":     info.Requests,
		"healthchecks": info.Healthchecks,
		"warns":        info.Warns,
		"errors":       info.Errors,
	}

	return result, nil
}

func extractTunnelID(fullURL, serverDomain string) string {
	url := strings.TrimPrefix(fullURL, "http://")
	url = strings.TrimPrefix(url, "https://")

	url = strings.TrimSuffix(url, "."+serverDomain)
	url = strings.TrimPrefix(url, serverDomain+"/")

	return url
}
