package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pedroborgesdev/tunnerse-cli/internal/server/debug"
)

type LogDetail struct {
	Key   string
	Value interface{}
}

var (
	logFiles = make(map[string]*os.File)
	logMutex sync.Mutex
)


func SetTunnelLogFile(tunnelID, logsDir string) error {
	logMutex.Lock()
	defer logMutex.Unlock()


	if _, exists := logFiles[tunnelID]; exists {
		return nil
	}


	if logsDir == "" {
		logsDir = "logs"
	}


	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}


	logPath := filepath.Join(logsDir, fmt.Sprintf("%s.log", tunnelID))
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logFiles[tunnelID] = file


	fmt.Printf("✓ Log file created: %s\n", logPath)

	return nil
}


func CloseTunnelLogFile(tunnelID string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if file, exists := logFiles[tunnelID]; exists {
		file.Sync()
		file.Close()
		delete(logFiles, tunnelID)
		fmt.Printf("✓ Log file closed: %s\n", tunnelID)
	}
}

func Log(level string, message string, details []LogDetail) {
	if !debug.DebugConfig.Debug && level == "DEBUG" {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	color := getLevelColor(level)
	reset := "\033[0m"


	var consoleMsg string
	if level == "DEBUG" {

		consoleMsg = fmt.Sprintf("%s %s%s:%d%s \n↳ %s%s%s - %s\n",
			timestamp, color, file, line, reset,
			color, level, reset, message)
	} else {

		consoleMsg = fmt.Sprintf("%s \n↳ %s%s%s - %s\n",
			timestamp, color, level, reset, message)
	}


	detailsMsg := ""
	for _, detail := range details {
		detailLine := fmt.Sprintf("  ↳ %s%s%s: %v\n", color, detail.Key, reset, detail.Value)
		consoleMsg += detailLine

		detailsMsg += detailLine
	}


	fmt.Print(consoleMsg)


	isTunnelLoop := strings.Contains(file, "tunnel_loop.go") || strings.Contains(file, "healthcheck.go")
	if isTunnelLoop {

		tunnelID := ""
		for _, detail := range details {
			if detail.Key == "tunnel_id" || detail.Key == "ID" {
				tunnelID = fmt.Sprintf("%v", detail.Value)
				break
			}
		}

		if tunnelID != "" {
			writeToLogFile(tunnelID, timestamp, level, file, line, message, detailsMsg)
		}
	}
}


func writeToLogFile(tunnelID, timestamp, level, file string, line int, message, details string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	logFile, exists := logFiles[tunnelID]
	fmt.Println(logFile, exists, tunnelID)
	if !exists {

		if err := SetTunnelLogFile(tunnelID, "/home/pedroborgezs/.tunnerse/logs/"); err != nil {
			return
		}
		logFile = logFiles[tunnelID]
	}

	if logFile != nil {
		color := getLevelColor(level)
		reset := "\033[0m"


		var fileMsg string
		if level == "DEBUG" {

			fileMsg = fmt.Sprintf("%s %s%s:%d%s \n↳ %s%s%s - %s\n%s",
				timestamp, color, file, line, reset, color, level, reset, message, details)
		} else {

			fileMsg = fmt.Sprintf("%s \n↳ %s%s%s - %s\n%s",
				timestamp, color, level, reset, message, details)
		}
		logFile.WriteString(fileMsg)
	}
}

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
	default:
		return "\033[35m"
	}
}
