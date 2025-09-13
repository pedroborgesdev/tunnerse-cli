package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"tunnerse/config"
	"tunnerse/database"
	"tunnerse/logger"
	"tunnerse/utils"
)

type HealthJob struct {
	repo *database.ActionsRepository
}

// NewHealthJob creates and returns a new instance of HealthJob.
func NewHealthJob(repo *database.ActionsRepository) *HealthJob {
	return &HealthJob{
		repo: repo,
	}
}

// StartHealthCheck launches a goroutine that periodically checks the health of the local API.
// It logs warnings on failures, and exits the program if failures exceed maxFails.
func (h *HealthJob) StartHealthCheck() {
	go func() {
		time.Sleep(5 * time.Second)

		failCount := 0
		maxFails := 10

		for {
			resp, err := http.Get(config.GetAddressURL())
			if err != nil {
				failCount++
				if isConnectionRefused(err) {
					logger.Log("WARN", "local API connection refused", []logger.LogDetail{
						{Key: "attempt", Value: toStr(failCount)},
					}, false)
					h.repo.UpdateWarnCount(config.GetTunnelID())
				} else {
					logger.Log("WARN", "health check failed", []logger.LogDetail{
						{Key: "attempt", Value: toStr(failCount)},
						{Key: "error", Value: err.Error()},
					}, false)
					h.repo.UpdateWarnCount(config.GetTunnelID())
				}

				if failCount >= maxFails {
					logger.Log("FATAL", "local API failed "+toStr(maxFails)+" times. closing tunnel.", nil, false)
					err := CloseConnection()
					if err != nil {
						logger.Log("FATAL", "error to close tunnel", nil, false)
					}

					println()
					os.Exit(0)
				}
			} else {
				resp.Body.Close()
				if failCount > 0 {
					logger.Log("INFO", "local API reestablished", []logger.LogDetail{}, false)
				}
				failCount = 0
			}

			time.Sleep(60 * time.Second)
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)

		for {
			h.PingToServer()
			time.Sleep(3550 * time.Second)
		}
	}()
}

// isConnectionRefused checks if the error string contains "connection refused".
func isConnectionRefused(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}

// toStr converts an integer to its string representation.
func toStr(n int) string {
	return fmt.Sprintf("%d", n)
}

// PingToServer sends a null request to persist the tunnel lifetime
func (h *HealthJob) PingToServer() {
	req, err := http.NewRequest("HEAD", utils.GetUrl("ping"), nil)
	if err != nil {
		logger.Log("ERROR", "error during send healthcheck challenge", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
		h.repo.UpdateErrorCount(config.GetTunnelID())
		return
	}

	req.Header.Set("Tunnerse", "healtcheck-question")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log("ERROR", "error during process healthcheck challenge", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
		h.repo.UpdateErrorCount(config.GetTunnelID())
		return
	}
	defer resp.Body.Close()

	if resp.Header.Get("Tunnerse") == "healthcheck-conclued" {
		logger.Log("HEALTHCHECK", "challenge has been overcome", []logger.LogDetail{}, false)
		h.repo.UpdateHealthcheckCount(config.GetTunnelID())
	} else {
		logger.Log("ERROR", "healthcheck challenge has failed", []logger.LogDetail{}, false)
		h.repo.UpdateErrorCount(config.GetTunnelID())
	}
}

// CloseConnection sends a request to the server to close the current tunnel connection.
func CloseConnection() error {
	payload := map[string]string{"name": config.GetTunnelID()}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	resp, err := http.Post(config.GetTunnelHTTPURL()+"/close", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("post register: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
