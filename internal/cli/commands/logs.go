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

	logPath := filepath.Join(config.GetLogsDir(), fmt.Sprintf("%s.log", tunnelID))

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logger.Log("FATAL", "Log file not found", []logger.LogDetail{
			{Key: "Tunnel_id", Value: tunnelID},
			{Key: "Path", Value: logPath},
			{Key: "Hint", Value: "Make sure the tunnel exists and is running"},
		}, false)
	}

	logger.Log("SUCCESS", fmt.Sprintf("Reading logs from tunnel '%s'", tunnelID), []logger.LogDetail{}, false)
	logger.Log("WARN", "Press Ctrl+C to stop", []logger.LogDetail{}, false)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go readLogFile(logPath)

	<-sigChan
	logger.Log("SUCCESS", "Stopped reading logs", []logger.LogDetail{}, false)
}

func readLogFile(logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		logger.Log("FATAL", "Failed to open log file", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Path", Value: logPath},
		}, false)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

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
		fmt.Print(line)
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return
		}
		fmt.Print(line)
	}
}
