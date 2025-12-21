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

package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/malczuuu/failbook/internal/metrics"
)

func LoggingAndMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		routePath := c.FullPath()
		if routePath == "" {
			routePath = path
		}

		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		log.Info().
			Str("method", method).
			Str("path", routePath).
			Interface("params", params).
			Interface("query", c.Request.URL.Query()).
			Int("status", status).
			Dur("latency", latency).
			Msg("processed http request")

		metrics.HTTPRequestsTotal.WithLabelValues(method, routePath, strconv.Itoa(status)).Inc()
		metrics.HTTPRequestDurationSeconds.WithLabelValues(method, routePath).Observe(latency.Seconds())
	}
}
