package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/dto"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/logger"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/utils"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/validators"

	"github.com/spf13/cobra"
)

// newTunnel representa o comando "new", que cria um túnel persistente.
var newTunnel = &cobra.Command{
	Use:                "new <tunnel_name> <local_port>",
	Short:              "Create a permanent tunnel connection (runs in background automatically)",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		validateNewArgs(args)
		startNewTunnel(args)
	},
}

// startNewTunnel registra o túnel via API local e salva no banco de dados local.
func startNewTunnel(args []string) {
	fmt.Printf(dto.Start)

	tunnelID := args[0]
	port := args[1]
	serverURL := "https://tunnerse.com" // Agora com protocolo explícito

	// Cria o payload para a API
	payload := map[string]string{
		"name":       tunnelID,
		"port":       port,
		"server_url": serverURL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Log("FATAL", "Failed to create request payload", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}

	// Faz a requisição POST para a API local
	apiURL := "http://localhost:9988/new"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		if utils.IsConnRefused(err) {
			logger.Log("FATAL", "Tunnerse local server is not online", []logger.LogDetail{
				{Key: "Hint", Value: "Make sure tunnerse-server is running and accessible on http://localhost:9988"},
			}, false)
		} else {
			logger.Log("FATAL", "Failed to connect to local API", []logger.LogDetail{
				{Key: "Error", Value: err.Error()},
			}, false)
		}
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
			Message   string `json:"message"`
			Subdomain bool   `json:"subdomain"`
			Tunnel    string `json:"tunnel"`
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
		logger.Log("FATAL", "Server returned error", []logger.LogDetail{
			{Key: "Code", Value: apiResponse.Code},
			{Key: "Message", Value: apiResponse.Message},
		}, false)
	}

	registeredTunnelID := apiResponse.Data.Tunnel
	isSubdomain := apiResponse.Data.Subdomain

	// Remove o protocolo do serverURL para construir a URL do túnel
	serverDomain := strings.TrimPrefix(serverURL, "http://")
	serverDomain = strings.TrimPrefix(serverDomain, "https://")

	// Extrai o protocolo do serverURL
	protocol := "http://"
	if strings.HasPrefix(serverURL, "https://") {
		protocol = "https://"
	}

	var tunnelURL string
	if isSubdomain {
		tunnelURL = fmt.Sprintf("%s%s.%s", protocol, registeredTunnelID, serverDomain)
	} else {
		tunnelURL = fmt.Sprintf("%s%s/%s", protocol, serverDomain, registeredTunnelID)
	}

	logger.Log("SUCCESS", "Tunnel is now running on server", []logger.LogDetail{
		{Key: "Tunnel_id", Value: registeredTunnelID},
		{Key: "Url", Value: tunnelURL},
	}, false)

	logger.Log("SUCCESS", "Tunnel is now managed by the server", []logger.LogDetail{}, false)
	logger.Log("INFO", "To see tunnel status, use 'tunnerse list'", []logger.LogDetail{}, false)
}

func validateNewArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateExposeArgs(args[0], args[1]); err != nil {
		logger.Log("ERROR", "Invalid args", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
		restoreTerminalAndExit(1)
	}
}

func restoreTerminalAndExit(code int) {
	utils.EnableInput()
	os.Exit(code)
}
