package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface that defines logging behavior
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

// zapLogger implements the Logger interface using zap
type zapLogger struct {
	logger *zap.Logger
}

var (
	// globalLogger is the singleton instance used by package-level functions
	globalLogger Logger
)

// Debug logs a debug level message
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info level message
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning level message
func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error level message
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal level message and then calls os.Exit(1)
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger with the provided fields
func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
	}
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// NewZapLogger creates a new Logger implemented with zap
func NewZapLogger(logLevel string, appEnvironment string) (Logger, error) {
	var zapConfig zap.Config
	isDevelopment := strings.ToLower(appEnvironment) == "staging"

	if isDevelopment {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Encoding = "json"
		zapConfig.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	level, err := zapcore.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid log level '%s', defaulting to 'info'. Error: %v\n", logLevel, err)
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapConfig.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize zap logger: %w", err)
	}

	return &zapLogger{logger: logger}, nil
}

// Init initializes the global logger
func Init(logLevel string, appEnvironment string) error {
	logger, err := NewZapLogger(logLevel, appEnvironment)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// L returns the global logger instance
func L() Logger {
	if globalLogger == nil {
		fmt.Fprintln(os.Stderr, "Logger not initialized. Call logger.Init() first. Using a default emergency logger.")
		_ = Init("info", "development")
	}
	return globalLogger
}

// Sync flushes any buffered log entries from the global logger
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}
