package logger

import (
	"fmt"
	"os"
	"time"
	"tunnerse/config"
)

// LogDetail represents a key-value pair for additional logging information.
type LogDetail struct {
	Key   string
	Value interface{}
}

// Log prints a formatted log message to the console with color, timestamp, level, message, and details.
func Log(level string, message string, details []LogDetail, showTime bool) {
	timestamp := time.Now().Format("2006/01/02-15:04:05")
	color := getLevelColor(level)
	reset := "\033[0m"

	if config.GetEnvBackgroundRunning() || config.GetEnvApplicationRunning() {
		if showTime {
			fmt.Printf("timestamp:[%s]level:[%s]info:[%s]", timestamp, level, message)
		} else {
			fmt.Printf("level:[%s]info:[%s]", level, message)
		}

		for _, detail := range details {
			fmt.Printf("detail:[%s=%s]", detail.Key, detail.Value)
		}
	} else {
		if showTime {
			fmt.Printf("%s %s%s \n↳ %s%s%s - %s",
				timestamp, color, reset,
				color, level, reset, message)
		} else {
			fmt.Printf("%s%s%s - %s",
				color, level, reset, message)
		}

		for _, detail := range details {
			fmt.Printf("\n%s%s%s: %v", color, detail.Key, reset, detail.Value)
		}
	}

	if level == "FATAL" {
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println()
}

// getLevelColor returns the ANSI color code corresponding to the log level.
func getLevelColor(level string) string {
	switch level {
	case "DEBUG":
		return "\033[36m"
	case "INFO":
		return "\033[32m"
	case "SUCCESS":
		return "\033[32m"
	case "WARN":
		return "\033[33m"
	case "HEALTHCHECK":
		return "\033[38;2;255;105;180m"
	case "ERROR":
		return "\033[31m"
	case "FATAL":
		return "\033[31m"
	default:
		return "\033[35m"
	}
}
