package commands

import (
	"fmt"
	"os"
	"strconv"
	"time"

	command_utils "tunnerse/commands/command_utils"
	"tunnerse/config"
	"tunnerse/dto"
	"tunnerse/logger"
	"tunnerse/models"
	"tunnerse/validators"

	"github.com/spf13/cobra"
)

// newTunnel representa o comando "new", que cria um túnel persistente.
var newTunnel = &cobra.Command{
	Use:                "new <tunnel_name> <local_port>",
	Short:              "Create a permanent tunnel connection (runs in background automatically)",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !config.GetEnvBackgroundRunning() {
			if config.GetEnvApplicationRunning() && len(args) > 1 {
				for i := 0; i < len(args)-1; i++ {
					args[i] = args[i+1]
				}
				args = args[:len(args)-1]
			}
		}

		config.SetTunnelID(args[0])
		config.SetAddressURL(args[1])

		config.SetServerURL("tunnerse.com")

		if config.GetEnvBackgroundRunning() {
			setupEnvironment()
			Server.StartTunnelLoop()
			return
		}

		validateNewArgs(args)
		startNewTunnel(args)
	},
}

// setupEnvironment define configurações vindas de variáveis de ambiente.
func setupEnvironment() {
	config.SetTunnelID(config.GetEnvTunneID())

	subdomain, _ := strconv.ParseBool(os.Getenv(config.GetEnvSubdomain()))
	config.SetSubdomainBool(subdomain)
}

// startNewTunnel registra e inicia o túnel em segundo plano.
func startNewTunnel(args []string) {
	fmt.Printf(dto.Start)

	tunnelID, isSubdomain, err := Server.RegisterTunnel()
	if err != nil {
		logger.Log("FATAL", "failed to register tunnel", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
	}

	config.SetTunnelID(tunnelID)
	config.SetSubdomainBool(isSubdomain)

	tunnelURL := command_utils.BuildTunnelURL()
	logger.Log("SUCCESS", "tunnel has been registered", []logger.LogDetail{
		{Key: "url", Value: tunnelURL},
	}, false)

	startOpts := &command_utils.StartOptions{
		LogDir:  config.GetExecPath() + "/" + tunnelID,
		LogName: "output",
	}

	pid, err := command_utils.StartInBackground(args, startOpts)
	if err != nil {
		logger.Log("FATAL", "failed to start tunnel in background", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
	}

	logger.Log("SUCCESS", "process has been initialized in background", []logger.LogDetail{
		{Key: "PID", Value: pid},
	}, false)

	tunnel := &models.Tunnel{
		ID:        tunnelID,
		Port:      args[1],
		Url:       tunnelURL,
		Domain:    config.GetServerURL(),
		Active:    true,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	info := &models.Info{
		ID:           tunnelID,
		Pid:          pid,
		Requests:     0,
		Healthchecks: 0,
		Warns:        0,
		Errors:       0,
	}

	if err := Repo.Create(tunnel, info); err != nil {
		logger.Log("ERROR", "tunnel not saved", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
	}

	fmt.Println("To see tunnel logs, use 'tunnerse logs <tunnel_id>'")
}

func validateNewArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateExposeArgs(args[0], args[1]); err != nil {
		logger.Log("ERROR", "invalid args", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
			{Key: "error2", Value: err.Error()},
		}, false)
		os.Exit(1)
	}
}
