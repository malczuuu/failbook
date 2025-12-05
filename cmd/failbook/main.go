package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/malczuuu/failbook/internal/config"
	"github.com/malczuuu/failbook/internal/logging"
	"github.com/malczuuu/failbook/internal/metrics"
	"github.com/malczuuu/failbook/internal/middleware"
)

func main() {
	cfg := config.Load()
	logging.ConfigureLogger(&cfg)

	metrics.Init()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.ZerologRecovery())
	router.Use(middleware.LoggingAndMetricsMiddleware())

	router.LoadHTMLGlob("./templates/*")

	if cfg.PrometheusEnabled {
		router.GET("/manage/prometheus", gin.WrapH(promhttp.Handler()))
		log.Info().Str("path", "/manage/prometheus").Msg("prometheus endpoint exposed")
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "Posts"})
	})
	router.GET("/problem", func(c *gin.Context) {
		c.HTML(http.StatusOK, "problem.tmpl", gin.H{"title": "Posts"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
	})

	router.NoMethod(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
	})

	addr := ":" + cfg.Port

	log.Info().Str("addr", addr).Msg("starting server")
	if err := router.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server exited")
	}
}
