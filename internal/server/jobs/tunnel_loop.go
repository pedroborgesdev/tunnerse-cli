package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pedroborgesdev/tunnerse-cli/internal/server/config"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/database"

	"github.com/pedroborgesdev/tunnerse-cli/internal/server/logger"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/models"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/repositories"
	"github.com/pedroborgesdev/tunnerse-cli/internal/server/utils"
)

type LoopJob struct {
	repo        *repositories.TunnelRepository
	ID          string
	tunnelURL   string
	localAPIURL string
	isSubdomain bool // true if this tunnel uses subdomain, false if uses path-based routing
	isQuick     bool
	stopChan    chan struct{}
	stopped     bool
	stopMu      sync.Mutex
}

// Stop para o tunnel loop e o healthcheck
func (s *LoopJob) Stop() {
	s.stopMu.Lock()
	defer s.stopMu.Unlock()

	if !s.stopped {
		close(s.stopChan)
		s.stopped = true
	}
}

func NewLoopJob(db *database.Database, ID string, port string, isSubdomain bool, serverDomain string, tunnelURL string, isQuick bool) *LoopJob {
	repo := repositories.NewTunnelRepository(db)

	// Se não for quick, busca a URL do túnel do banco de dados
	var finalTunnelURL string
	if !isQuick {
		tunnel, err := repo.GetTunnel(ID)
		if err != nil {
			logger.Log("ERROR", "failed to get tunnel from database", []logger.LogDetail{
				{Key: "tunnel_id", Value: ID},
				{Key: "error", Value: err.Error()},
			})
			return nil
		}
		finalTunnelURL = tunnel.Url
	} else {
		// Para quick, usa a URL passada como parâmetro
		finalTunnelURL = tunnelURL
	}

	// Constrói a URL da API local diretamente da porta
	localAPIURL := fmt.Sprintf("http://localhost:%s", port)

	job := &LoopJob{
		repo:        repo,
		ID:          ID,
		tunnelURL:   finalTunnelURL,
		localAPIURL: localAPIURL,
		isSubdomain: isSubdomain, // Store whether this specific tunnel uses subdomain
		isQuick:     isQuick,
		stopChan:    make(chan struct{}),
	}

	return job
}

func printDecodedBody(resp *models.ResponseData) {
	fmt.Println(string(resp.Body))
}

func (s *LoopJob) SendResponseToServer(data *models.ResponseData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	responseURL := s.tunnelURL + "/response"
	// printDecodedBody(data)
	logger.Log("DEBUG", "sending response to server", []logger.LogDetail{
		{Key: "tunnel_id", Value: s.ID},
		{Key: "response_url", Value: responseURL},
	})

	_, err = http.Post(responseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log("ERROR", "failed to send response", []logger.LogDetail{
			{Key: "tunnel_id", Value: s.ID},
			{Key: "response_url", Value: responseURL},
			{Key: "error", Value: err.Error()},
		})
	}
	return err
}

// StartTunnelLoop starts the main loop that continuously fetches, processes, and responds to tunnel requests.
func (s *LoopJob) StartTunnelLoop() {
	if err := logger.SetTunnelLogFile(s.ID, config.LogsDir); err != nil {
		logger.Log("ERROR", "failed to create log file", []logger.LogDetail{
			{Key: "tunnel_id", Value: s.ID},
			{Key: "logs_dir", Value: config.LogsDir},
			{Key: "error", Value: err.Error()},
		})
	}
	defer logger.CloseTunnelLogFile(s.ID)

	// Garante que os mapas serão limpos quando o loop terminar
	defer func() {
		delete(config.QuickTunnelURLs, s.ID)
		config.RemoveActiveJob(s.ID)
	}()

	logger.Log("INFO", "starting tunnel loop", []logger.LogDetail{
		{Key: "tunnel_id", Value: s.ID},
	})

	var errorTimestamps []time.Time
	const rateLimitCount = 10
	const rateLimitWindow = 10 * time.Second

	go s.healthcheckLocalAPI()
	go s.pingToServer()

	for {
		select {
		case <-s.stopChan:
			logger.Log("INFO", "tunnel loop stopped by external signal", []logger.LogDetail{
				{Key: "tunnel_id", Value: s.ID},
			})
			return
		default:
			// Continue com o fluxo normal
		}

		now := time.Now()
		errorTimestamps = filterRecent(errorTimestamps, now, rateLimitWindow)

		if len(errorTimestamps) >= rateLimitCount {
			logger.Log("FATAL", "target server did not respond, the tunnel was closed", []logger.LogDetail{
				{Key: "tunnel_id", Value: s.ID},
			})
			s.Stop()
			break
		}

		reqData, err := s.FetchRequest()

		if reqData != nil && err != nil && err.Error() == "healthcheck-question" {
			respData := &models.ResponseData{
				StatusCode: 204,
				Headers: map[string][]string{
					"Tunnerse": {"healthcheck-conclued"},
				},
				Body:  nil,
				Token: reqData.Token,
			}
			err = s.SendResponseToServer(respData)
			if err != nil {
				logger.Log("ERROR", "failed to send healthcheck response", []logger.LogDetail{
					{Key: "tunnel_id", Value: s.ID},
					{Key: "error", Value: err.Error()},
				})
			}
			continue
		}

		if reqData == nil && err == nil {
			continue
		}

		if err != nil {
			if err.Error() == "tunnel has closed by server" {
				logger.Log("FATAL", "tunnel has closed by server", []logger.LogDetail{
					{Key: "tunnel_id", Value: s.ID},
				})
				break
			}
			if err.Error() == "response-time-exceeded" {
				logger.Log("FATAL", "reponse time exceeded, the tunnel was closed", []logger.LogDetail{
					{Key: "tunnel_id", Value: s.ID},
				})
				break
			}
			errorTimestamps = append(errorTimestamps, now)
			continue
		}

		respData, err := s.ForwardToLocal(reqData)
		if err != nil {
			logger.Log("WARN", "failed to forward request to local API", []logger.LogDetail{
				{Key: "tunnel_id", Value: s.ID},
				{Key: "error", Value: err.Error()},
			})

			// Envia resposta de erro ao servidor para não deixar a requisição pendurada
			errorResp := &models.ResponseData{
				StatusCode: http.StatusServiceUnavailable,
				Headers: map[string][]string{
					"Content-Type": {"text/plain; charset=utf-8"},
					"Tunnerse":     {"local-api-error"},
				},
				Token: reqData.Token,
			}

			sendErr := s.SendResponseToServer(errorResp)
			if sendErr != nil {
				logger.Log("ERROR", "failed to send error response", []logger.LogDetail{
					{Key: "error", Value: sendErr.Error()},
				})
			}
			continue
		}

		if !s.isQuick {
			s.repo.UpdateRequestCount(s.ID)
		}

		err = s.SendResponseToServer(respData)
		if err != nil {
			logger.Log("FATAL", "error during send response to server", []logger.LogDetail{
				{Key: "tunnel_id", Value: s.ID},
				{Key: "error", Value: err.Error()},
			})
			break
		}
	}
}

// FetchRequest fetches the incoming request data from the tunnel server.
func (s *LoopJob) FetchRequest() (*models.RequestData, error) {
	fetchURL := s.tunnelURL + "/tunnel"
	logger.Log("DEBUG", "fetching request from server", []logger.LogDetail{
		{Key: "tunnel_id", Value: s.ID},
		{Key: "fetch_url", Value: fetchURL},
	})

	resp, err := http.Get(fetchURL)
	if err != nil {
		logger.Log("ERROR", "failed to fetch request", []logger.LogDetail{
			{Key: "tunnel_id", Value: s.ID},
			{Key: "fetch_url", Value: fetchURL},
			{Key: "error", Value: err.Error()},
		})
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
		switch value[0] {
		case "healthcheck-question":
			return &requestData, fmt.Errorf("healthcheck-question")
		case "tunnel-not-found":
			return nil, fmt.Errorf("notfound")
		case "tunnel-timeout":
			return nil, fmt.Errorf("timeout")
		case "tunnel-working":
			return nil, fmt.Errorf("working")
		}
	}

	fmt.Println(requestData.Path)

	return &requestData, nil
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func (s *LoopJob) ForwardToLocal(req *models.RequestData) (*models.ResponseData, error) {
	if isTunnerseDemoPath(req.Path) {
		demoResp, err := serveDemoHTML(req.Path)
		if err != nil {
			return nil, err
		}
		demoResp.Token = req.Token
		return demoResp, nil
	}

	path := req.Path
	tunnelPrefix := "/" + s.ID + "/"
	if strings.HasPrefix(path, tunnelPrefix) {
		path = "/" + strings.TrimPrefix(path, tunnelPrefix)
	}
	url := fmt.Sprintf("%s%s", s.localAPIURL, path)

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

	if !s.isSubdomain {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			tunnelName := s.ID
			// body = utils.InjectBaseHref(body, tunnelName)
			body = utils.RewriteAbsolutePaths(body, tunnelName)
		}
	}

	headers := make(map[string][]string)
	for key, values := range resp.Header {
		if strings.ToLower(key) == "content-length" {
			continue
		}
		headers[key] = values
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = []string{"text/html; charset=utf-8"}
	}

	var respData *models.ResponseData
	respData = &models.ResponseData{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
		Token:      req.Token,
	}

	return respData, nil
}

func isTunnerseDemoPath(path string) bool {
	p := path
	if p == "" {
		return false
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	if p == "/tunnerse" || strings.HasPrefix(p, "/tunnerse/") {
		return true
	}
	return false
}

func serveDemoHTML(requestPath string) (*models.ResponseData, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	demoPath := filepath.Join(home, ".tunnerse", "static", "demo.html")
	data, err := os.ReadFile(demoPath)
	if err != nil {
		return nil, err
	}

	headers := map[string][]string{
		"Content-Type": {"text/html; charset=utf-8"},
		"Tunnerse":     {"demo"},
	}

	return &models.ResponseData{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       data,
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

func (s *LoopJob) closeConnection() error {
	payload := map[string]string{"name": s.ID}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	resp, err := http.Post(s.tunnelURL+"/close", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("post register: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
