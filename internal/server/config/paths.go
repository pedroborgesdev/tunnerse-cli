package config

import (
	"os"
	"path/filepath"
)



func GetUserDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {

		return ".tunnerse"
	}
	return filepath.Join(homeDir, ".tunnerse")
}


func GetLogsDir() string {
	return filepath.Join(GetUserDataDir(), "logs")
}


func GetDatabasePath() string {
	return filepath.Join(GetUserDataDir(), "db.sqlite")
}


func EnsureDataDirExists() error {
	dataDir := GetUserDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	logsDir := GetLogsDir()
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return err
	}

	return nil
}
