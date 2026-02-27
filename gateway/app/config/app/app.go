package app

import (
	"context"
	"gateway/app/core/gateway"
	"gateway/app/dataprovider/lib"
	"gateway/app/dataprovider/telemetry"
	"net/http"
	"time"
)

type Config struct {
	ServerAddr   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	IsDevelopment bool

	OtelEnabled        bool
	OtelEndpoint       string
	OtelServiceName    string
	OtelServiceVersion string
	OtelDebugExporter  bool
}

type App struct {
	config           Config
	httpServer       *http.ServeMux
	logger           gateway.LoggerGateway
	isRunning        bool
	otelSDKShutdown  func(ctx context.Context) error
	otelLogsShutdown func(ctx context.Context) error
}

func NewApp(config Config) *App {
	return &App{
		config:    config,
		isRunning: false,
	}
}

func (a *App) Run() error {
	if a.config.IsDevelopment || a.config.OtelEnabled {
		shutdownLogs, err := setupOTelLogs(
			context.Background(),
			a.config.OtelEndpoint,
			a.config.OtelServiceName,
			a.config.OtelServiceVersion,
		)
		if err != nil {
			println("Failed to initialize OTEL logs:", err.Error())
		} else {
			a.otelLogsShutdown = shutdownLogs
		}
	}

	logger, err := lib.NewZapLogger(a.config.IsDevelopment)
	if err != nil {
		return err
	}
	gateway.SetLogger(logger)
	a.logger = logger

	a.logger.Info("Starting server...")

	if a.config.OtelEnabled {
		a.logger.Info("Starting OpenTelemetry...",
			"endpoint", a.config.OtelEndpoint,
			"serviceName", a.config.OtelServiceName,
			"serviceVersion", a.config.OtelServiceVersion,
		)
		shutdownSDK, err := telemetry.SetupTracingAndMetrics(context.Background(), telemetry.Config{
			Endpoint:       a.config.OtelEndpoint,
			Insecure:       true,
			ServiceName:    a.config.OtelServiceName,
			ServiceVersion: a.config.OtelServiceVersion,
		})
		if err != nil {
			a.logger.Error("Failed to start OTEL, continuing without starting", "error", err)
		} else {
			a.otelSDKShutdown = shutdownSDK
		}

		if a.otelLogsShutdown != nil {
			a.logger.Info("OpenTelemetry enabled", "endpoint", a.config.OtelEndpoint)
		}
	} else {
		a.logger.Info("OpenTelemetry unabled")
		if a.otelLogsShutdown != nil {
			a.logger.Info("OpenTelemetry enabled [DEV_MODE]", "endpoint", a.config.OtelEndpoint)
		}
	}

	mux := http.NewServeMux()

	a.httpServer = mux

	a.logger.Info("Server started")
	return nil
}

func (a *App) ServeHttp() error {
	a.isRunning = true

	return http.ListenAndServe(a.config.ServerAddr, a.httpServer)
}

func (a *App) Shutdown(ctx context.Context) error {
	if !a.isRunning {
		a.logger.Debug("Shutdown called but server is not running")
		return nil
	}

	a.logger.Info("Shutting down server...")

	if a.otelSDKShutdown != nil {
		a.logger.Info("Shutting down OTEL")
		if err := a.otelSDKShutdown(ctx); err != nil {
			a.logger.Error("Failed to shutdown OTEL", "error", err)
		}
	}

	if a.otelLogsShutdown != nil {
		a.logger.Info("Shutting down OTEL Logs")
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := a.otelLogsShutdown(ctx); err != nil {
			a.logger.Error("Failed to shutdown OTEL Logs", "error", err)
		}
	}

	a.isRunning = false
	a.logger.Info("Server shutdown with success")
	return nil
}
