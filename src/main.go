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
	fmt.Print(dto.Start)

	validator := validators.NewArgsValidator()
	server := servers.NewServerService()
	keyboard := jobs.NewKeyboardJob()

	err := validator.ValidateUsageArgs(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	tunnelID, localAddress := os.Args[1], os.Args[2]

	config.SetTunnelID(tunnelID)
	config.SetAddressURL(localAddress)
	config.SetServerURL("tunnerse.com")

	keyboard.CloseKeyboardJob()

	err = validator.ValidateExposeArgs(tunnelID, localAddress)
	if err != nil {
		logger.LogError("FATAL", err, true)
	}

	_, _, err = server.RegisterTunnel()
	if err != nil {
		logger.LogError("REGISTER TUNNEL", err, true)
		return
	}

	var tunnelURL string
	if config.GetSubdomainBool() {
		tunnelURL = fmt.Sprintf("------> http://%s.%s", config.GetTunnelID(), config.GetServerURL())
	} else {
		fmt.Print(dto.BetaWarn)
		tunnelURL = fmt.Sprintf("------> http://%s/%s", config.GetServerURL(), config.GetTunnelID())
	}

	logger.Log("INFO", "tunnel has been registered", []logger.LogDetail{{Key: "url", Value: tunnelURL}})

	server.StartTunnelLoop()
}
