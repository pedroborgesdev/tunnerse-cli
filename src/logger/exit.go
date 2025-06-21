package logger

import (
	"fmt"
	"os"
)

// LogError logs an error message with context and exits the program if shouldExit is true.
func LogError(context string, err error, shouldExit bool) {
	if err == nil {
		return
	}

	level := "ERROR"
	if shouldExit {
		level = "FATAL"
	}

	Log(level, "error in "+context, []LogDetail{
		{Key: "error", Value: err.Error()},
	})

	if shouldExit {
		fmt.Println()
		os.Exit(1)
	}
}
