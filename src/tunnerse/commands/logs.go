package commands

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"tunnerse/config"
	"tunnerse/jobs"
	"tunnerse/logger"
	"tunnerse/validators"

	"github.com/spf13/cobra"
)

// newTunnel representa o comando "new", que cria um túnel persistente.
var logsTunnel = &cobra.Command{
	Use:                "logs <tunnel_id>",
	Short:              "show tunnel logs in real time",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		jobs.CloseKeyboardJob()
		config.SetTunnelID(args[0])
		validateLogsArgs(args)
		logsRun()
	},
}

// validateArgs verifica se os argumentos fornecidos são válidos.
func validateLogsArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateTunnelID(args[0]); err != nil {
		fmt.Printf("invalid args: %s\n", err.Error())
		os.Exit(1)
	}
}

func logsRun() {
	path := fmt.Sprint(config.GetExecPath() + "/" + config.GetTunnelID() + "/" + "output.log")

	file, err := os.Open(path)
	if err != nil {
		logger.Log("FATAL", "failed to read log file", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		fmt.Print(line)
	}
}
