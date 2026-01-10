package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/jobs"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/logger"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/validators"

	"github.com/spf13/cobra"
)

// delTunnel representa o comando "del", que deleta um túnel inativo.
var delTunnel = &cobra.Command{
	Use:                "del <tunnel_id>",
	Short:              "delete an inactive tunnel from database",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobs.CloseKeyboardJob()
		validateDelArgs(args)
		delRun(args[0])
	},
}

// validateDelArgs verifica se os argumentos fornecidos são válidos.
func validateDelArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateTunnelID(args[0]); err != nil {
		logger.Log("FATAL", "Invalid arguments", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}
}

func delRun(tunnelID string) {
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

	// Faz a requisição DELETE para a API local
	apiURL := "http://localhost:9988/delete"
	req, err := http.NewRequest("DELETE", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log("FATAL", "Failed to create request", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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
			Message  string `json:"message"`
			TunnelID string `json:"tunnel_id"`
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
		// Verifica se é túnel ativo
		if apiResponse.Message == "Invalid request data" &&
			strings.Contains(fmt.Sprintf("%v", apiResponse.Data), "still active") {
			logger.Log("ERROR", "Tunnel is still active", []logger.LogDetail{
				{Key: "Tunnel_id", Value: tunnelID},
				{Key: "Hint", Value: "Use 'tunnerse kill " + tunnelID + "' first"},
			}, false)
			return
		}
		// Outros erros
		logger.Log("FATAL", "Server returned error", []logger.LogDetail{
			{Key: "Code", Value: apiResponse.Code},
			{Key: "Message", Value: apiResponse.Message},
		}, false)
	}

	logger.Log("SUCCESS", "Tunnel has been deleted from database", []logger.LogDetail{
		{Key: "Tunnel_id", Value: tunnelID},
	}, false)
}
