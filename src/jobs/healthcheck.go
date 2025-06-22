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
	"tunnerse/logger"
	"tunnerse/utils"
)

type HealthJob struct {
	urls *utils.UrlsUtils
}

// NewHealthJob creates and returns a new instance of HealthJob.
func NewHealthJob() *HealthJob {
	return &HealthJob{
		urls: utils.NewUrlsUtils(),
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
					})
				} else {
					logger.Log("WARN", "health check failed", []logger.LogDetail{
						{Key: "attempt", Value: toStr(failCount)},
						{Key: "error", Value: err.Error()},
					})
				}

				if failCount >= maxFails {
					logger.Log("FATAL", "local API failed "+toStr(maxFails)+" times. closing tunnel.", nil)
					err := h.CloseConnection()
					if err != nil {
						logger.Log("FATAL", "error to close tunnel", nil)
					}

					println()
					os.Exit(0)
				}
			} else {
				resp.Body.Close()
				if failCount > 0 {
					logger.Log("INFO", "local API reestablished", []logger.LogDetail{})
				}
				failCount = 0
			}

			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			h.PingToServer()
			time.Sleep(3400 * time.Second)
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
	req, err := http.NewRequest("HEAD", h.urls.GetUrl("ping"), nil)
	if err != nil {
		logger.LogError("FATAL", err, true)
		return
	}

	req.Header.Set("Tunnerse", "healtcheck-question")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.LogError("FATAL", err, true)
		return
	}
	defer resp.Body.Close()

	if resp.Header.Get("Tunnerse") == "healthcheck-conclued" {
		logger.Log("INFO", "the server health check challenge has been overcome", []logger.LogDetail{
			{Key: "status", Value: resp.Header.Get("Tunnerse")},
		})
	}
}

// CloseConnection sends a request to the server to close the current tunnel connection.
func (h *HealthJob) CloseConnection() error {
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
