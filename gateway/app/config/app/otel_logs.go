package app

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func setupOTelLogs(ctx context.Context, endpoint, serviceName, serviceVersion string) (func(context.Context) error, error) {
	var shutdowns []func(context.Context) error

	shutdown := func(ctx context.Context) error {
		var joined error
		for _, fn := range shutdowns {
			joined = errors.Join(joined, fn(ctx))
		}
		return joined
	}

	exp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return shutdown, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)
	if err != nil {
		return shutdown, err
	}

	processor := log.NewBatchProcessor(
		exp,
		log.WithMaxQueueSize(2048),
		log.WithExportInterval(1*time.Second),
	)

	lp := log.NewLoggerProvider(
		log.WithProcessor(processor),
		log.WithResource(res),
	)
	shutdowns = append(shutdowns, lp.Shutdown, exp.Shutdown)

	global.SetLoggerProvider(lp)

	return shutdown, nil
}
