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
	Version           string
}

func Load() Config {
	return Config{
		Port:              getenv("FAILBOOK_PORT", "12001"),
		LogLevel:          getenv("FAILBOOK_LOG_LEVEL", "info"),
		HealthEnabled:     getenv("FAILBOOK_HEALTH_ENABLED", "false") == "true",
		PrometheusEnabled: getenv("FAILBOOK_PROMETHEUS_ENABLED", "false") == "true",
		ProblemsDir:       getenv("FAILBOOK_PROBLEM_DOCS_DIR", "./problem-docs"),
		BaseHref:          getenv("FAILBOOK_BASE_HREF", "/"),
		Version:           getenv("FAILBOOK_VERSION", "unspecified"),
	}
}

func getenv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
