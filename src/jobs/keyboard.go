package jobs

import (
	"os"
	"os/signal"
	"syscall"
	"tunnerse/logger"
)

type KeyboardJob struct {
	healtcheck *HealthJob
}

// NewKeyboardJob creates and returns a new instance of KeyboardJob with an embedded HealthJob.
func NewKeyboardJob() *KeyboardJob {
	return &KeyboardJob{
		healtcheck: NewHealthJob(),
	}
}

// CloseKeyboardJob listens for OS interrupt signals (SIGINT, SIGTERM) and gracefully closes the tunnel before exiting.
func (k *KeyboardJob) CloseKeyboardJob() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		println()

		logger.Log("WARN", "closing tunnerse", nil)

		err := k.healtcheck.CloseConnection()
		if err != nil {
			return
		}

		println()
		os.Exit(0)
	}()
}
