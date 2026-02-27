package main

import (
	"context"
	"gateway/app/config/app"
	"gateway/app/config/env"
	"gateway/app/core/gateway"
	"gateway/app/dataprovider/lib"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		os.Exit(1)
	}

	logger, err := lib.NewZapLogger(env.GetBool("IS_DEVELOPMENT", true))
	if err != nil {
		os.Exit(1)
	}
	gateway.SetLogger(logger)

	appConfig := app.Config{
		ServerAddr:   env.GetString("SERVER_ADDR", ":8080"),
		ReadTimeout:  env.GetDuration("SERVER_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: env.GetDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  env.GetDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),

		IsDevelopment: env.GetBool("IS_DEVELOPMENT", true),

		OtelEnabled:        env.GetBool("OTEL_ENABLED", false),
		OtelEndpoint:       env.GetString("OTEL_ENDPOINT", "localhost:4317"),
		OtelServiceName:    env.GetString("OTEL_SERVICE_NAME", "gateway"),
		OtelServiceVersion: env.GetString("OTEL_SERVICE_VERSION", "0.1.0"),
		OtelDebugExporter:  env.GetBool("OTEL_DEBUG_EXPORTER", false),
	}
	application := app.NewApp(appConfig)
	if err := application.Run(); err != nil {
		logger.Error("Failed to start application", "error", err)
		os.Exit(1)
	}

	serverError := make(chan error, 1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverError:
		logger.Error("Failed to start server", "error", err)
	case sig := <-quit:
		logger.Info("Received shutdown signal", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown application", "error", err)
		os.Exit(1)
	}

	logger.Info("Application shutdown")
}
