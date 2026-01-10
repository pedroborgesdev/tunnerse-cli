package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/config"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/jobs"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/logger"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/validators"

	"github.com/spf13/cobra"
)

// logsTunnel representa o comando "logs", que mostra os logs do túnel em tempo real.
var logsTunnel = &cobra.Command{
	Use:                "logs <tunnel_id>",
	Short:              "show tunnel logs in real time",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobs.CloseKeyboardJob()
		validateLogsArgs(args)
		logsRun(args[0])
	},
}

// validateLogsArgs verifica se os argumentos fornecidos são válidos.
func validateLogsArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateTunnelID(args[0]); err != nil {
		logger.Log("FATAL", "Invalid arguments", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}
}

func logsRun(tunnelID string) {
	logger.Log("INFO", "Reading tunnel logs...", []logger.LogDetail{}, false)

	// Define o caminho do arquivo de log usando o diretório do usuário
	logPath := filepath.Join(config.GetLogsDir(), fmt.Sprintf("%s.log", tunnelID))

	// Verifica se o arquivo existe
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logger.Log("FATAL", "Log file not found", []logger.LogDetail{
			{Key: "Tunnel_id", Value: tunnelID},
			{Key: "Path", Value: logPath},
			{Key: "Hint", Value: "Make sure the tunnel exists and is running"},
		}, false)
	}

	logger.Log("SUCCESS", fmt.Sprintf("Reading logs from tunnel '%s'", tunnelID), []logger.LogDetail{}, false)
	logger.Log("WARN", "Press Ctrl+C to stop", []logger.LogDetail{}, false)

	// Configurar handler para Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Inicia tail -f do arquivo de log
	go readLogFile(logPath)

	// Aguarda sinal de interrupção
	<-sigChan
	logger.Log("SUCCESS", "Stopped reading logs", []logger.LogDetail{}, false)
}

// readLogFile lê um arquivo de log do início e continua monitorando (tail -f)
func readLogFile(logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		logger.Log("FATAL", "Failed to open log file", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Path", Value: logPath},
		}, false)
	}
	defer file.Close()

	// Lê o arquivo do início (para mostrar histórico)
	reader := bufio.NewReader(file)

	// Lê tudo que já existe
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Log("FATAL", "Failed to read log file", []logger.LogDetail{
				{Key: "Error", Value: err.Error()},
			}, false)
		}
		// Imprime a linha (que já vem com cores ANSI do logger)
		fmt.Print(line)
	}

	// Agora continua lendo novas linhas (tail -f)
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
