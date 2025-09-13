package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"tunnerse-server/config"
	"tunnerse-server/database"
	"tunnerse-server/jobs"
	"tunnerse-server/models"
	"tunnerse-server/repositories"
	"tunnerse-server/utils"
)

type TunnelService struct {
	healtcheck *jobs.HealthJob
	repo       *repositories.TunnelRepository
}

func NewTunnelService(db *database.Database) *TunnelService {
	repo := repositories.NewTunnelRepository(db)
	return &TunnelService{
		healtcheck: jobs.NewHealthJob(repo),
		repo:       repo,
	}
}

// CloseConnection sends a request to close the current tunnel session on the server.
func (s *TunnelService) CloseConnection() error {
	payload := map[string]string{"name": config.GetTunnelID()}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	resp, err := http.Post(config.GetTunnelHTTPSURL()+"/close", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("post register: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// RegisterTunnel sends a registration request to the server and returns the assigned tunnel ID.
func (s *TunnelService) RegisterTunnel(name, port, server_domain string) (string, bool, error) {
	payload := map[string]string{"name": name}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", false, fmt.Errorf("encode JSON: %w", err)
	}

	resp, err := http.Post(utils.GetUrl("register", name, false, server_domain), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", false, fmt.Errorf("post register: %w", err)
	}
	defer resp.Body.Close()

	var result models.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", false, fmt.Errorf("decode register response. probably tunnerse server is offline: %w", err)
	}

	config.LocalRoutes[result.Data.Tunnel] = fmt.Sprintf("http://localhost:%s", port)
	tunnel_url := utils.BuildTunnelURL(result.Data.Tunnel, server_domain, result.Data.Subdomain)

	tunnel := &models.Tunnel{
		ID:        result.Data.Tunnel,
		Port:      port,
		Url:       tunnel_url,
		Domain:    server_domain,
		Active:    true,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	info := &models.Info{
		ID:           result.Data.Tunnel,
		Requests:     0,
		Healthchecks: 0,
		Warns:        0,
		Errors:       0,
	}

	if err := s.repo.Create(tunnel, info); err != nil {
		return "", false, fmt.Errorf("tunnel not saved: %w", err)
	}

	go func() {
		jobs.NewLoopJob(s.repo.DB, result.Data.Tunnel, result.Data.Subdomain, server_domain).StartTunnelLoop()
	}()

	return tunnel_url, result.Data.Subdomain, nil
}
