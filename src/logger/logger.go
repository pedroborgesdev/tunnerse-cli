package logger

import (
	"fmt"
	"time"
)

// LogDetail represents a key-value pair for additional logging information.
type LogDetail struct {
	Key   string
	Value interface{}
}

// Log prints a formatted log message to the console with color, timestamp, level, message, and details.
func Log(level string, message string, details []LogDetail) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	color := getLevelColor(level)
	reset := "\033[0m"

	fmt.Printf("%s %s%s \n↳ %s%s%s - %s\n",
		timestamp, color, reset,
		color, level, reset, message)

	for _, detail := range details {
		fmt.Printf("  %s%s%s: %v\n", color, detail.Key, reset, detail.Value)
	}
}

// getLevelColor returns the ANSI color code corresponding to the log level.
func getLevelColor(level string) string {
	switch level {
	case "DEBUG":
		return "\033[36m"
	case "INFO":
		return "\033[32m"
	case "WARN":
		return "\033[33m"
	case "ERROR":
		return "\033[31m"
	case "FATAL":
		return "\033[31m"
	default:
		return "\033[35m"
	}
}
