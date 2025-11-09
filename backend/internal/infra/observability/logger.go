package observability

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/trustmedis/mini-attendance/internal/config"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Sync() error
}

type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func InitLogger() Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	logger, _ := cfg.Build()
	return &ZapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func InitLoggerWithConfig(cfg config.LoggingConfig) Logger {
	var zapCfg zap.Config

	if cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Set log level
	switch cfg.Level {
	case "debug":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, _ := zapCfg.Build()
	return &ZapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	l.sugar.Infow(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	l.sugar.Errorw(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...interface{}) {
	l.sugar.Warnw(msg, fields...)
}

func (l *ZapLogger) Debug(msg string, fields ...interface{}) {
	l.sugar.Debugw(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...interface{}) {
	l.sugar.Fatalw(msg, fields...)
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

// StructuredLogger for Loki integration
type StructuredLogger struct {
	baseLogger Logger
	labels     map[string]string
}

func NewStructuredLogger(cfg config.LoggingConfig) *StructuredLogger {
	return &StructuredLogger{
		baseLogger: InitLoggerWithConfig(cfg),
		labels:     cfg.LokiLabels,
	}
}

func (sl *StructuredLogger) LogWithLabels(level string, msg string, fields map[string]interface{}) {
	// Add labels to fields
	for k, v := range sl.labels {
		fields[fmt.Sprintf("label_%s", k)] = v
	}

	// Convert map to slice for zap
	var zapFields []interface{}
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}

	switch level {
	case "info":
		sl.baseLogger.Info(msg, zapFields...)
	case "error":
		sl.baseLogger.Error(msg, zapFields...)
	case "warn":
		sl.baseLogger.Warn(msg, zapFields...)
	case "debug":
		sl.baseLogger.Debug(msg, zapFields...)
	}
}
