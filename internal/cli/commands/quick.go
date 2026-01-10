package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/config"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/dto"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/logger"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/utils"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/validators"

	"github.com/spf13/cobra"
)

// quickTunnel representa o comando "quick", que inicia o túnel diretamente no terminal atual.
var quickTunnel = &cobra.Command{
	Use:   "quick <tunnel_name> <local_port>",
	Short: "Start a quick tunnel on current terminal (no database)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		startQuickTunnel(args)
	},
}

type QuickResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Tunnel    string `json:"tunnel"`
		Subdomain bool   `json:"subdomain"`
	} `json:"data"`
	Status int `json:"status"`
}

// startQuickTunnel executa o fluxo do túnel rápido, validando e registrando via API.
func startQuickTunnel(args []string) {
	utils.Clear()

	fmt.Printf(dto.Welcome)
	fmt.Printf(dto.Start)

	validateQuickArgs(args)

	tunnelName := args[0]
	port := args[1]
	serverURL := "https://tunnerse.com"

	// Faz requisição para a API local
	payload := map[string]string{
		"name":       tunnelName,
		"port":       port,
		"server_url": serverURL,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Log("FATAL", "Failed to create request payload", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}

	resp, err := http.Post("http://localhost:9988/quick", "application/json", bytes.NewBuffer(data))
	if err != nil {
		logger.Log("FATAL", "Failed to connect to server", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Hint", Value: "Make sure tunnerse-server is running"},
		}, false)
	}
	defer resp.Body.Close()

	var result QuickResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Log("FATAL", "Failed to decode server response", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}

	if result.Status != 200 {
		logger.Log("FATAL", "Server returned error", []logger.LogDetail{
			{Key: "Message", Value: result.Message},
			{Key: "Code", Value: result.Code},
		}, false)
	}

	tunnelID := result.Data.Tunnel
	isSubdomain := result.Data.Subdomain

	// Constrói a URL do túnel usando o mesmo protocolo e domínio do server_url
	// Remove o protocolo para extrair apenas o domínio
	serverDomain := strings.TrimPrefix(serverURL, "http://")
	serverDomain = strings.TrimPrefix(serverDomain, "https://")

	// Detecta o protocolo do serverURL
	protocol := "http://"
	if strings.HasPrefix(serverURL, "https://") {
		protocol = "https://"
	}

	var tunnelURL string
	if isSubdomain {
		tunnelURL = fmt.Sprintf("%s%s.%s", protocol, tunnelID, serverDomain)
	} else {
		tunnelURL = fmt.Sprintf("%s%s/%s", protocol, serverDomain, tunnelID)
	}

	logger.Log("SUCCESS", "Quick tunnel created successfully!", []logger.LogDetail{
		{Key: "Tunnel URL", Value: tunnelURL},
	}, false)
	logger.Log("WARN", "Press Ctrl+C to stop", []logger.LogDetail{}, false)

	// Aguarda um pouco para o arquivo de log ser criado
	time.Sleep(500 * time.Millisecond)

	// Define o caminho do arquivo de log usando o diretório do usuário
	logPath := filepath.Join(config.GetLogsDir(), fmt.Sprintf("%s.log", tunnelID))

	// Aguarda o arquivo ser criado
	maxWait := 5 * time.Second
	startWait := time.Now()
	for {
		if _, err := os.Stat(logPath); err == nil {
			break
		}
		if time.Since(startWait) > maxWait {
			logger.Log("WARN", "Log file not found", []logger.LogDetail{
				{Key: "Path", Value: logPath},
			}, false)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Configurar handler para Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Inicia tail -f do arquivo de log
	go tailLogFile(logPath)

	// Aguarda sinal de interrupção
	<-sigChan

	fmt.Println()
	logger.Log("INFO", "Stopping tunnel...", []logger.LogDetail{}, false)
	stopTunnel(tunnelID)
	logger.Log("SUCCESS", "Quick tunnel stopped", []logger.LogDetail{}, false)

	restoreTerminalAndExit(1)
}

// stopTunnel envia requisição para matar o túnel
func stopTunnel(tunnelID string) {
	payload := map[string]string{
		"tunnel_id": tunnelID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Log("ERROR", "Failed to create kill request", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
		return
	}

	resp, err := http.Post("http://localhost:9988/kill", "application/json", bytes.NewBuffer(data))
	if err != nil {
		logger.Log("ERROR", "Failed to stop tunnel", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log("ERROR", "Failed to read response", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Log("ERROR", "Failed to parse response", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
		return
	}

	if resp.StatusCode != 200 {
		logger.Log("ERROR", "Failed to stop tunnel", []logger.LogDetail{
			{Key: "Error", Value: fmt.Sprintf("%v", result["error"])},
		}, false)
		return
	}
}

// tailLogFile faz tail -f de um arquivo de log
func tailLogFile(logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		logger.Log("ERROR", "Failed to open log file", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Path", Value: logPath},
		}, false)
		return
	}
	defer file.Close()

	// Vai para o final do arquivo
	file.Seek(0, io.SeekEnd)

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Aguarda mais conteúdo
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return
		}
		// Imprime a linha (que já vem com cores ANSI do logger)
		fmt.Print(line)
	}
}

// validateQuickArgs verifica se os argumentos fornecidos são válidos.
func validateQuickArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateExposeArgs(args[0], args[1]); err != nil {
		logger.Log("FATAL", "Invalid arguments", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}
}

func restoreTerminalAndExitQuick(code int) {
	utils.EnableInput()
	os.Exit(code)
}
