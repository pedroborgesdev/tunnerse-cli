package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"tunnerse/config"
	"tunnerse/jobs"
	"tunnerse/logger"
	"tunnerse/models"
	"tunnerse/utils"
)

type ServerService struct {
	healtcheck *jobs.HealthJob
	rewrite    *utils.RewriteUtils
}

// NewServerService creates and returns a new instance of ServerService with health check initialized.
func NewServerService() *ServerService {
	return &ServerService{
		healtcheck: jobs.NewHealthJob(),
		rewrite:    utils.NewRewriteUtils(),
	}
}

// GetUrl returns the appropriate server URL based on the requested method type.
func (s *ServerService) GetUrl(method string) string {
	switch method {
	case "register":
		return "http://" + config.GetServerURL() + "/register"
	case "response":
		return config.GetTunnelHTTPSURL() + "/response"
	case "fetch":
		return config.GetTunnelHTTPSURL() + "/tunnel"
	}
	return "undefined"
}

// CloseConnection sends a request to close the current tunnel session on the server.
func (s *ServerService) CloseConnection() error {
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
func (s *ServerService) RegisterTunnel() (string, bool, error) {
	payload := map[string]string{"name": config.GetTunnelID()}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", false, fmt.Errorf("encode JSON: %w", err)
	}

	resp, err := http.Post(s.GetUrl("register"), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", false, fmt.Errorf("post register: %w", err)
	}
	defer resp.Body.Close()

	var result models.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", false, fmt.Errorf("decode register response. probably tunnerse server is offline")
	}

	config.SetTunnelID(result.Data.Tunnel)
	config.SetSubdomainBool(result.Data.Subdomain)

	return result.Data.Tunnel, result.Data.Subdomain, nil
}

// SendResponseToServer sends the local response back to the server after processing the request.
func (s *ServerService) SendResponseToServer(data *models.ResponseData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = http.Post(s.GetUrl("response"), "application/json", bytes.NewBuffer(jsonData))
	return err
}

// StartTunnelLoop starts the main loop that continuously fetches, processes, and responds to tunnel requests.
func (s *ServerService) StartTunnelLoop() {
	var errorTimestamps []time.Time
	const rateLimitCount = 5
	const rateLimitWindow = 10 * time.Second

	s.healtcheck.StartHealthCheck()

	for {
		now := time.Now()
		filtered := errorTimestamps[:0]
		for _, t := range errorTimestamps {
			if now.Sub(t) <= rateLimitWindow {
				filtered = append(filtered, t)
			}
		}
		errorTimestamps = filtered

		if len(errorTimestamps) >= rateLimitCount {
			logger.LogError("RECEIVE REQUEST LIMIT", fmt.Errorf("max attempt exceeded"), true)
		}

		reqData, err := s.FetchRequest()
		if err != nil {
			if err.Error() == "tunnel has closed by server" {
				logger.LogError("TUNNEL CONNECTION", err, true)
			}

			errorTimestamps = append(errorTimestamps, time.Now())
			logger.LogError("RECEIVE REQUEST", err, false)
			continue
		}

		respData, err := s.ForwardToLocal(reqData)
		if err != nil {
			logger.Log("WARN", "failed to FOWARD REQUEST", []logger.LogDetail{
				{Key: "error", Value: err.Error()},
			})
			continue
		}

		err = s.SendResponseToServer(respData)
		if err != nil {
			logger.LogError("SEND RESPONSE TO SERVER", err, false)
			continue
		}
	}
}

// FetchRequest fetches the incoming request data from the tunnel server.
func (s *ServerService) FetchRequest() (*models.RequestData, error) {
	resp, err := http.Get(s.GetUrl("fetch"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(bodyBytes, []byte("<!-- Tunnerse by @pedroborgezs")) {
		return nil, fmt.Errorf("tunnel has closed by server")
	}

	var requestData models.RequestData
	err = json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		return nil, fmt.Errorf("unexpected response by server %s", err.Error())
	}
	return &requestData, nil
}

// ForwardToLocal forwards the fetched request to the local service and captures the response.
func (s *ServerService) ForwardToLocal(req *models.RequestData) (*models.ResponseData, error) {
	client := &http.Client{}

	url := fmt.Sprintf(config.GetAddressURL() + "/" + req.Path)

	request, err := http.NewRequest(req.Method, url, bytes.NewBuffer([]byte(req.Body)))
	if err != nil {
		return nil, err
	}

	for key, values := range req.Headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		tunnelName := config.GetTunnelID()

		body = s.rewrite.InjectBaseHref(body, tunnelName)
		body = s.rewrite.RewriteAbsolutePaths(body, tunnelName)
	}

	return &models.ResponseData{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}
