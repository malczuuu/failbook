package config

import (
	"os"
)

type Config struct {
	Port     string
	LogLevel string
}

func Load() Config {
	port := getenv("FAILBOOK_PORT", "12001")
	logLevel := getenv("FAILBOOK_LOG_LEVEL", "info")

	return Config{
		Port:     port,
		LogLevel: logLevel,
	}
}

func getenv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
