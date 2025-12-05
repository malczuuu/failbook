package config

import (
	"os"
)

type Config struct {
	Port              string
	LogLevel          string
	PrometheusEnabled bool
}

func Load() Config {
	return Config{
		Port:              getenv("FAILBOOK_PORT", "12001"),
		LogLevel:          getenv("FAILBOOK_LOG_LEVEL", "info"),
		PrometheusEnabled: getenv("FAILBOOK_PROMETHEUS_ENABLED", "false") == "true",
	}
}

func getenv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
