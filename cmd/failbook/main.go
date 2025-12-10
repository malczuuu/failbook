package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/malczuuu/failbook/internal/config"
	"github.com/malczuuu/failbook/internal/health"
	"github.com/malczuuu/failbook/internal/logging"
	"github.com/malczuuu/failbook/internal/markdown"
	"github.com/malczuuu/failbook/internal/metrics"
	"github.com/malczuuu/failbook/internal/middleware"
	"github.com/malczuuu/failbook/internal/problems"
)

func main() {
	cfg := config.Load()
	logging.ConfigureLogger(&cfg)

	problemRegistry, err := problems.LoadFromDirectory(cfg.ProblemsDir)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load error configurations")
	}

	metrics.Init()

	healthStatus := health.NewStatus()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.ZerologRecovery())
	router.Use(middleware.LoggingAndMetricsMiddleware())

	router.LoadHTMLGlob("./templates/*")

	if cfg.HealthEnabled {
		router.GET("/manage/health/live", health.LivenessHandler())
		log.Info().Str("path", "/manage/health/live").Msg("liveness endpoint exposed")

		router.GET("/manage/health/ready", health.ReadinessHandler(healthStatus))
		log.Info().Str("path", "/manage/health/ready").Msg("readiness endpoint exposed")
	}

	if cfg.PrometheusEnabled {
		router.GET("/manage/prometheus", gin.WrapH(promhttp.Handler()))
		log.Info().Str("path", "/manage/prometheus").Msg("prometheus endpoint exposed")
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "API Error Documentation",
			"problems": problemRegistry.GetAll(),
			"baseHref": cfg.BaseHref,
		})
	})

	router.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		problem, exists := problemRegistry.Get(id)
		if !exists {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
				"baseHref": cfg.BaseHref,
			})
			return
		}

		c.HTML(http.StatusOK, "problem.tmpl", gin.H{
			"problem":         problem,
			"baseHref":        cfg.BaseHref,
			"descriptionHTML": markdown.RenderToHTML(problem.Description),
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"baseHref": cfg.BaseHref,
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"baseHref": cfg.BaseHref,
		})
	})

	addr := ":" + cfg.Port

	srv := &http.Server{Addr: addr, Handler: router}
	healthStatus.SetReady()

	go func() {
		log.Info().Str("addr", addr).Msg("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server exited with error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Info().Str("signal", sig.String()).Msg("commencing graceful shutdown")

	healthStatus.SetNotReady()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("graceful shutdown completed")
}
