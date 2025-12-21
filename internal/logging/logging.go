// Copyright (c) 2025 Damian Malczewski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// SPDX-License-Identifier: MIT

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
