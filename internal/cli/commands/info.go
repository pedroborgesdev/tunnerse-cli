package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/jobs"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/logger"

	"github.com/spf13/cobra"
)

// infoTunnel representa o comando "info", que exibe informações do túnel.
var infoTunnel = &cobra.Command{
	Use:                "info <tunnel_id>",
	Short:              "show tunnel information",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobs.CloseKeyboardJob()
		infoRun(args[0])
	},
}

func infoRun(tunnelID string) {
	// Cria o payload para a API
	payload := map[string]string{
		"tunnel_id": tunnelID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Log("FATAL", "Failed to create request payload", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}

	// Faz a requisição POST para a API local
	apiURL := "http://localhost:9988/info"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log("FATAL", "Failed to connect to local server", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Hint", Value: "Make sure tunnerse-server is running"},
		}, false)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log("FATAL", "Failed to read server response", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}

	// Parse da resposta
	var apiResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Info struct {
				ID           string `json:"id"`
				Port         string `json:"port"`
				Url          string `json:"url"`
				Domain       string `json:"domain"`
				Active       bool   `json:"active"`
				CreatedAt    string `json:"created_at"`
				Requests     int    `json:"requests"`
				Healthchecks int    `json:"healthchecks"`
				Warns        int    `json:"warns"`
				Errors       int    `json:"errors"`
			} `json:"info"`
		} `json:"data"`
		Status int `json:"status"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Log("FATAL", "Failed to parse server response", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Response", Value: string(body)},
		}, false)
	}

	// Verifica se houve erro na API
	if apiResponse.Code != "success" {
		// Verifica se é erro de túnel não encontrado
		if apiResponse.Code == "not_found" {
			logger.Log("ERROR", "Tunnel not found", []logger.LogDetail{
				{Key: "Tunnel_id", Value: tunnelID},
				{Key: "Hint", Value: "Use 'tunnerse list' to see available tunnels"},
			}, false)
			return
		}
		// Outros erros
		logger.Log("FATAL", "Server returned error", []logger.LogDetail{
			{Key: "Code", Value: apiResponse.Code},
			{Key: "Message", Value: apiResponse.Message},
		}, false)
	}

	info := apiResponse.Data.Info

	status := "Inactive"
	if info.Active {
		status = "Active"
	}

	fmt.Printf(
		"\033[36mID:           \033[0m%s\n"+
			"\033[36mPort:         \033[0m%s\n"+
			"\033[36mURL:          \033[0m%s\n"+
			"\033[36mDomain:       \033[0m%s\n"+
			"\033[36mStatus:       \033[0m%s\n"+
			"\033[36mCreatedAt:    \033[0m%s\n\n"+
			"\033[32mRequests:     \033[0m%v\n"+
			"\033[38;2;255;105;180mHealthchecks: \033[0m%v\n"+
			"\033[33mWarns:        \033[0m%v\n"+
			"\033[31mErrors:       \033[0m%v\n",
		info.ID, info.Port, info.Url, info.Domain, status, info.CreatedAt,
		info.Requests, info.Healthchecks, info.Warns, info.Errors,
	)
}
