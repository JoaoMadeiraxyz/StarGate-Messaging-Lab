package lib

import (
	"gateway/app/core/gateway"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(isDevelopment bool) (gateway.LoggerGateway, error) {
	var baseLogger *zap.Logger
	var err error

	if isDevelopment {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		baseLogger, err = cfg.Build()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		baseLogger, err = cfg.Build()
	}
	if err != nil {
		return nil, err
	}

	baseCore := baseLogger.Core()
	otelCore := otelzap.NewCore("gateway/app")
	teeCore := zapcore.NewTee(baseCore, otelCore)

	zLogger := zap.New(teeCore)

	_, _ = zap.RedirectStdLogAt(zLogger, zapcore.InfoLevel)

	sugar := zLogger.Sugar()
	return &ZapLogger{logger: sugar}, nil
}

func (z *ZapLogger) Debug(msg string, fields ...interface{}) {
	z.logger.Debugw(msg, fields...)
}

func (z *ZapLogger) Info(msg string, fields ...interface{}) {
	z.logger.Infow(msg, fields...)
}

func (z *ZapLogger) Warn(msg string, fields ...interface{}) {
	z.logger.Warnw(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...interface{}) {
	z.logger.Errorw(msg, fields...)
}

func (z *ZapLogger) Fatal(msg string, fields ...interface{}) {
	z.logger.Fatalw(msg, fields...)
}

func (z *ZapLogger) With(fields ...interface{}) gateway.LoggerGateway {
	return &ZapLogger{
		logger: z.logger.With(fields...),
	}
}
