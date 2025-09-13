package jobs

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tunnerse/logger"
	"tunnerse/utils"
)

// CloseKeyboardJob listens for OS interrupt signals (SIGINT, SIGTERM) and gracefully closes the tunnel before exiting.
func CloseKeyboardAndTunnelJob() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		logger.Log("WARN", "closing tunnel", nil, false)

		err := CloseConnection()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		println()
		utils.EnableInput()
		os.Exit(0)
	}()
}

func CloseKeyboardJob() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		println()
		utils.EnableInput()
		os.Exit(0)
	}()
}
