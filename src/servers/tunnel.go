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
	urls       *utils.UrlsUtils
}

// NewServerService creates and returns a new instance of ServerService with health check initialized.
func NewServerService() *ServerService {
	return &ServerService{
		healtcheck: jobs.NewHealthJob(),
		rewrite:    utils.NewRewriteUtils(),
		urls:       utils.NewUrlsUtils(),
	}
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
func (s *ServerService) RegisterTunnel() error {
	payload := map[string]string{"name": config.GetTunnelID()}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	resp, err := http.Post(s.urls.GetUrl("register"), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("post register: %w", err)
	}
	defer resp.Body.Close()

	var result models.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode register response. probably tunnerse server is offline")
	}

	config.SetTunnelID(result.Data.Tunnel)
	config.SetSubdomainBool(result.Data.Subdomain)

	return nil
}

// SendResponseToServer sends the local response back to the server after processing the request.
func (s *ServerService) SendResponseToServer(data *models.ResponseData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = http.Post(s.urls.GetUrl("response"), "application/json", bytes.NewBuffer(jsonData))
	return err
}

// StartTunnelLoop starts the main loop that continuously fetches, processes, and responds to tunnel requests.
func (s *ServerService) StartTunnelLoop() {
	var errorTimestamps []time.Time
	const rateLimitCount = 10
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

			errorTimestamps = append(errorTimestamps, time.Now())
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
			logger.Log("FATAL", "error during send response to server", []logger.LogDetail{})
			continue
		}
	}
}

// FetchRequest fetches the incoming request data from the tunnel server.
func (s *ServerService) FetchRequest() (*models.RequestData, error) {
	resp, err := http.Get(s.urls.GetUrl("fetch"))
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

	// if bytes.HasPrefix(bodyBytes, []byte("<!-- Tunnerse by @pedroborgezs")) {
	// 	return nil, fmt.Errorf("tunnel has closed by server")
	// }

	var requestData models.RequestData

	err = json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		fmt.Println(string(bodyBytes))
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
func (s *ServerService) ForwardToLocal(req *models.RequestData) (*models.ResponseData, error) {
	url := fmt.Sprintf("%s%s", config.GetAddressURL(), req.Path)

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

			body = s.rewrite.InjectBaseHref(body, tunnelName)
			body = s.rewrite.RewriteAbsolutePaths(body, tunnelName)
		}
	}

	return &models.ResponseData{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}
