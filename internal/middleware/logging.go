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

		log.Info().Str("method", method).Str("path", routePath).Int("status", status).Dur("latency", latency).Msg("request")

		metrics.HTTPRequestsTotal.WithLabelValues(method, routePath, strconv.Itoa(status)).Inc()
		metrics.HTTPRequestDurationSeconds.WithLabelValues(method, routePath).Observe(latency.Seconds())
	}
}
