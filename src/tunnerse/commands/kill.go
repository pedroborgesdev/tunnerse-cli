package commands

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"tunnerse/jobs"
	"tunnerse/logger"
	"tunnerse/validators"

	"github.com/spf13/cobra"
)

// newTunnel representa o comando "new", que cria um túnel persistente.
var killTunnel = &cobra.Command{
	Use:                "kill <tunnel_id>",
	Short:              "kill tunnel process",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobs.CloseKeyboardJob()
		validateKillArgs(args)
		killRun(args[0])
	},
}

// validateArgs verifica se os argumentos fornecidos são válidos.
func validateKillArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateTunnelID(args[0]); err != nil {
		fmt.Printf("invalid args: %s\n", err.Error())
		os.Exit(1)
	}
}

func killRun(tunnelID string) {
	logger.Log("INFO", "trying get tunnel informations", []logger.LogDetail{}, false)

	pid, err := Repo.GetPID(tunnelID)
	if err == nil && pid == 0 {
		logger.Log("FATAL", "tunnel not found", []logger.LogDetail{}, false)
	}
	if err != nil {
		logger.Log("FATAL", "failed to get tunnel PID", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
	}

	logger.Log("WARN", "finishing tunnel process", []logger.LogDetail{
		{Key: "PID", Value: pid},
	}, false)

	proc, err := os.FindProcess(pid)
	if err != nil {
		logger.Log("SUCCESS", "tunnel already inactive", []logger.LogDetail{}, false)
	} else {
		if runtime.GOOS == "windows" {
			err = proc.Kill()
		} else {
			err = proc.Signal(syscall.SIGKILL)
		}
		if err != nil {
			logger.Log("FATAL", "failed to finish tunnel process", []logger.LogDetail{
				{Key: "error", Value: err.Error()},
			}, false)
		}
	}

	err = Repo.UpdateTunnelStatus(tunnelID, false)
	if err != nil {
		logger.Log("FATAL", "failed to update tunnel status", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
	}

	logger.Log("SUCCESS", "tunnel updated to inactive", []logger.LogDetail{}, false)
}
