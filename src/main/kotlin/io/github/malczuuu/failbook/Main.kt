package io.github.malczuuu.failbook

import io.javalin.Javalin
import io.javalin.http.Context
import io.javalin.micrometer.MicrometerPlugin
import io.micrometer.prometheusmetrics.PrometheusConfig
import io.micrometer.prometheusmetrics.PrometheusMeterRegistry
import org.slf4j.Logger
import org.slf4j.LoggerFactory

val log: Logger = LoggerFactory.getLogger("io.github.malczuuu.problem.registry.Main")

val prometheus = PrometheusMeterRegistry(PrometheusConfig.DEFAULT)
val micrometer = MicrometerPlugin { it.registry = prometheus }

fun main() {
  val app =
      Javalin.create { config ->
        config.showJavalinBanner = false
        config.useVirtualThreads = true

        config.registerPlugin(micrometer)

        config.requestLogger.http { ctx, ms ->
          log.info(
              "Handled HTTP request status={}, method={}, requestPath={} in {}ms",
              ctx.status(),
              ctx.method(),
              ctx.path() + if (ctx.queryString() != null) "?${ctx.queryString()}" else "",
              ms,
          )
        }

        config.http.asyncTimeout = 10_000L
      }

  app.get("/") { it.result("Hello World!") }
  app.get("/manage/prometheus") { servePrometheus(it) }
  app.get("/manage/health") { serveHealth(it) }

  Runtime.getRuntime().addShutdownHook(Thread { app.stop() })
  app.start(7070)
}

fun servePrometheus(ctx: Context) {
  ctx.contentType("text/plain; version=0.0.4; charset=utf-8").result(prometheus.scrape())
}

fun serveHealth(ctx: Context) {
  ctx.contentType("application/health+json").json(Health())
}
