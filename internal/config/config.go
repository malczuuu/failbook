package config

import (
	"os"
)

type Config struct {
	Port              string
	LogLevel          string
	HealthEnabled     bool
	PrometheusEnabled bool
	ProblemsDir       string
	BaseHref          string
}

func Load() Config {
	return Config{
		Port:              getenv("FAILBOOK_PORT", "12001"),
		LogLevel:          getenv("FAILBOOK_LOG_LEVEL", "info"),
		HealthEnabled:     getenv("FAILBOOK_HEALTH_ENABLED", "false") == "true",
		PrometheusEnabled: getenv("FAILBOOK_PROMETHEUS_ENABLED", "false") == "true",
		ProblemsDir:       getenv("FAILBOOK_PROBLEM_DOCS_DIR", "./problem-docs"),
		BaseHref:          getenv("FAILBOOK_BASE_HREF", ""),
	}
}

func getenv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
