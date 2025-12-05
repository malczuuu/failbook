package logging

import (
	"os"
	"time"

	"github.com/malczuuu/failbook/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogger(cfg *config.Config) {
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = false

	logLevel := parseLevel(cfg.LogLevel)
	log.Logger = zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Logger()
}

func parseLevel(levelStr string) zerolog.Level {
	logLevel, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	return logLevel
}
