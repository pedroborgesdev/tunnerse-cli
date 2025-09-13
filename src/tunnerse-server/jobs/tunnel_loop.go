package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"tunnerse-server/config"
	"tunnerse-server/database"

	"tunnerse-server/logger"
	"tunnerse-server/models"
	"tunnerse-server/repositories"
	"tunnerse-server/utils"
)

type LoopJob struct {
	healthcheck  *HealthJob
	repo         *repositories.TunnelRepository
	ID           string
	isSubdomain  bool
	serverDomain string
}

func NewLoopJob(db *database.Database, ID string, isSubdomain bool, serverDomain string) *LoopJob {
	repo := repositories.NewTunnelRepository(db)
	return &LoopJob{
		healthcheck:  NewHealthJob(repo),
		repo:         repo,
		ID:           ID,
		isSubdomain:  isSubdomain,
		serverDomain: serverDomain,
	}
}

// SendResponseToServer sends the local response back to the server after processing the request.
func (s *LoopJob) SendResponseToServer(data *models.ResponseData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = http.Post(utils.GetUrl("response", s.ID, s.isSubdomain, s.serverDomain), "application/json", bytes.NewBuffer(jsonData))
	return err
}

// StartTunnelLoop starts the main loop that continuously fetches, processes, and responds to tunnel requests.
func (s *LoopJob) StartTunnelLoop() {
	logger.Log("INFO", "starting tunnel loop", []logger.LogDetail{})

	var errorTimestamps []time.Time
	const rateLimitCount = 10
	const rateLimitWindow = 10 * time.Second

	// s.healthcheck.StartHealthCheck()

	for {
		if len(config.LocalRoutes) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		now := time.Now()
		errorTimestamps = filterRecent(errorTimestamps, now, rateLimitWindow)

		if len(errorTimestamps) >= rateLimitCount {
			logger.Log("FATAL", "target server did not respond, the tunnel was closed", []logger.LogDetail{})
		}

		reqData, err := s.FetchRequest()
		if reqData == nil && err == nil {
			respData := models.ResponseData{
				StatusCode: 204,
				Headers: map[string][]string{
					"Tunnerse": {"healtcheck-response"},
				},
				Body: nil,
			}
			s.SendResponseToServer(&respData)
			continue
		}

		if err != nil {
			if err.Error() == "tunnel has closed by server" {
				logger.Log("FATAL", "tunnel has closed by server", []logger.LogDetail{})
			}
			if err.Error() == "response-time-exceeded" {
				logger.Log("FATAL", "reponse time exceeded, the tunnel was closed", []logger.LogDetail{})
			}
			errorTimestamps = append(errorTimestamps, now)
			continue
		}

		respData, err := s.ForwardToLocal(reqData)
		if err != nil {
			logger.Log("WARN", "failed to FOWARD REQUEST", []logger.LogDetail{})
			continue
		}

		s.repo.UpdateRequestCount(s.ID)

		err = s.SendResponseToServer(respData)
		if err != nil {
			logger.Log("FATAL", "error during send response to server", []logger.LogDetail{})
			continue
		}
	}
}

// FetchRequest fetches the incoming request data from the tunnel server.
func (s *LoopJob) FetchRequest() (*models.RequestData, error) {
	resp, err := http.Get(utils.GetUrl("fetch", s.ID, s.isSubdomain, s.serverDomain))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusGatewayTimeout:
			return nil, fmt.Errorf("response-time-exceeded")
		default:
			return nil, fmt.Errorf("unexpected response by server")
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var requestData models.RequestData

	err = json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		return nil, fmt.Errorf("unexpected response by server: %s", err.Error())
	}

	value, ok := requestData.Headers["Tunnerse"]
	if ok && len(value) > 0 {
		if value[0] == "healtcheck-question" {
			return nil, nil
		}

		switch value[0] {
		case "healthcheck-question":
			return nil, nil
		case "tunnel-not-found":
			return nil, fmt.Errorf("notfound")
		case "tunnel-timeout":
			return nil, fmt.Errorf("timeout")
		case "tunnel-working":
			return nil, fmt.Errorf("working")
		}
	}

	return &requestData, nil
}

var httpClient = &http.Client{}

// ForwardToLocal forwards the fetched request to the local service and captures the response.
func (s *LoopJob) ForwardToLocal(req *models.RequestData) (*models.ResponseData, error) {
	target, ok := config.LocalRoutes[s.ID]
	var local_url string

	if ok {
		local_url = target
	}

	url := fmt.Sprintf("%s%s", local_url, req.Path)

	request, err := http.NewRequest(req.Method, url, bytes.NewBuffer([]byte(req.Body)))
	if err != nil {
		return nil, err
	}

	for key, values := range req.Headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !config.GetSubdomainBool() {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			tunnelName := config.GetTunnelID()
			body = utils.InjectBaseHref(body, tunnelName)
			body = utils.RewriteAbsolutePaths(body, tunnelName)
		}
	}

	// Remove Content-Length para evitar conflito
	headers := make(map[string][]string)
	for key, values := range resp.Header {
		if strings.ToLower(key) == "content-length" {
			continue
		}
		headers[key] = values
	}

	// Garantir Content-Type se não existir
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = []string{"text/html; charset=utf-8"}
	}

	return &models.ResponseData{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
	}, nil
}

func filterRecent(timestamps []time.Time, now time.Time, window time.Duration) []time.Time {
	filtered := timestamps[:0]
	for _, t := range timestamps {
		if now.Sub(t) <= window {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
