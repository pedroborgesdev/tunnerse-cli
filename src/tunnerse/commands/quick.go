package commands

import (
	"fmt"
	"os"

	command_utils "tunnerse/commands/command_utils"
	"tunnerse/config"
	"tunnerse/dto"
	"tunnerse/jobs"
	"tunnerse/logger"
	"tunnerse/utils"
	"tunnerse/validators"

	"github.com/spf13/cobra"
)

// quickTunnel representa o comando "quick", que inicia o túnel diretamente no terminal atual.
var quickTunnel = &cobra.Command{
	Use:   "quick <tunnel_name> <local_port>",
	Short: "Start a quick tunnel on current terminal",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		startQuickTunnel(args)
	},
}

// startQuickTunnel executa o fluxo do túnel rápido, validando, registrando e iniciando.
func startQuickTunnel(args []string) {
	utils.Clear()

	fmt.Printf(dto.Welcome)
	fmt.Printf(dto.Start)

	validateQuickArgs(args)

	config.SetTunnelID(args[0])
	config.SetAddressURL(args[1])
	config.SetServerURL("tunnerse.com")

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

	jobs.CloseKeyboardAndTunnelJob()

	Server.StartTunnelLoop()
}

// validateArgs verifica se os argumentos fornecidos são válidos.
func validateQuickArgs(args []string) {
	validator := validators.NewArgsValidator()

	if err := validator.ValidateExposeArgs(args[0], args[1]); err != nil {
		fmt.Printf("invalid args: %s\n", err.Error())
		os.Exit(1)
	}
}
