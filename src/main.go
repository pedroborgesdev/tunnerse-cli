package main

import (
	"fmt"
	"os"
	"tunnerse/config"
	"tunnerse/dto"
	"tunnerse/jobs"
	"tunnerse/logger"
	"tunnerse/servers"
	"tunnerse/validators"
)

func main() {
	fmt.Print(dto.Welcome)

	validator := validators.NewArgsValidator()
	server := servers.NewServerService()
	keyboard := jobs.NewKeyboardJob()

	keyboard.CloseKeyboardJob()

	msg := validator.ValidateUsageArgs(os.Args)
	if msg != "" {
		fmt.Print(msg)
		os.Exit(0)
	}

	tunnelID, localAddress := os.Args[1], os.Args[2]

	err := validator.ValidateExposeArgs(tunnelID, localAddress)
	if err != nil {
		logger.Log("FATAL", "argument error", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		})
	}

	config.SetTunnelID(tunnelID)
	config.SetAddressURL(localAddress)
	config.SetServerURL("tunnerse.com")

	fmt.Print(dto.Start)

	_, _, err = server.RegisterTunnel()
	if err != nil {
		logger.Log("FATAL", "failed to register tunnel", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		})
		return
	}

	var tunnelURL string
	if config.GetSubdomainBool() {
		tunnelURL = fmt.Sprintf("------> http://%s.%s", config.GetTunnelID(), config.GetServerURL())
	} else {
		fmt.Print(dto.BetaWarn)
		tunnelURL = fmt.Sprintf("------> http://%s/%s", config.GetServerURL(), config.GetTunnelID())
	}

	logger.Log("SUCCESS", "tunnel has been registered", []logger.LogDetail{{Key: "url", Value: tunnelURL}})

	server.StartTunnelLoop()
}
