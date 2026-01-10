package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/pedroborgesdev/tunnerse-cli/internal/server/logger"
)

var (
	LogsDir string
)

type Config struct {
	HTTPPort string

	SUBDOMAIN     bool
	WARNS_ON_HTML bool

	TUNNEL_LIFE_TIME            int
	TUNNEL_INACTIVITY_LIFE_TIME int
}

var AppConfig Config

func LoadAppConfig() error {

	if err := EnsureDataDirExists(); err != nil {
		logger.Log("ERROR", "failed to create data directory", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		})
		return err
	}

	LogsDir = GetLogsDir()

	err := godotenv.Load()
	if err != nil {
		logger.Log("DEBUG", "Error on read .env file", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
	}

	AppConfig = Config{
		HTTPPort: getEnvStr("HTTPPort", "9988"),

		SUBDOMAIN:     getEnvBool("SUBDOMAIN", false),
		WARNS_ON_HTML: getEnvBool("WARNS_ON_HTML", true),

		TUNNEL_LIFE_TIME:            getEnvInt("TUNNEL_LIFE_TIME", 86400),
		TUNNEL_INACTIVITY_LIFE_TIME: getEnvInt("TUNNEL_INACTIVITY_LIFE_TIME", 86400),
	}

	logger.Log("ENV", "Defined environment variables", []logger.LogDetail{
		{Key: "HTTPPort", Value: AppConfig.HTTPPort},
		{Key: "SUBDOMAIN", Value: AppConfig.SUBDOMAIN},
	})
	return nil
}

func getEnvStr(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
