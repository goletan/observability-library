// /observability/logger/logger.go
package logger

import (
	"fmt"
	"sync"

	"github.com/goletan/config"
	"github.com/goletan/observability/types"
	"github.com/goletan/observability/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger   *zap.Logger
	once     sync.Once
	scrubber = utils.NewScrubber()
	cfg      *types.ObservabilityConfig
)

// InitLogger initializes the default logger. It will panic if the logger fails to initialize.
func InitLogger() {
	once.Do(func() {
		var err error
		zapConfig := zap.NewProductionConfig()

		cfg = &types.ObservabilityConfig{}
		err = config.LoadConfig("Observability", cfg, nil)
		if err != nil {
			fmt.Printf("Warning: failed to load observability configuration, using defaults: %v\n", err)
		}

		if cfg.Logger.Level != "" {
			if err := zapConfig.Level.UnmarshalText([]byte(cfg.Logger.Level)); err != nil {
				fmt.Printf("Invalid log level: %v, defaulting to INFO\n", cfg.Logger.Level)
				zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
			}
		}

		logger, err = zapConfig.Build()
		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}
	})
}

// SetLogger allows setting a custom logger, primarily used for testing.
func SetLogger(customLogger *zap.Logger) {
	once.Do(func() {
		// Ensure thread safety by using sync.Once to avoid race conditions.
		logger = customLogger
	})
	logger = customLogger
}

// ensureLoggerInitialized ensures that the logger is initialized before use.
func ensureLoggerInitialized() {
	if logger == nil {
		panic("Logger not initialized. Call InitLogger first.")
	}
}

// Info logs an info level message with additional fields.
func Info(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Info(message, fields...)
}

// Debug logs a debug level message with additional fields.
func Debug(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Debug(message, fields...)
}

// Error logs an error level message with additional fields.
func Error(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Error(message, fields...)
}

// Warn logs a warning level message with additional fields.
func Warn(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Warn(message, fields...)
}

// Fatal logs a fatal error and exits the application.
func Fatal(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Fatal(message, fields...)
}

// WithContext adds contextual fields to the log entries.
func WithContext(context map[string]interface{}) []zap.Field {
	ensureLoggerInitialized()
	fields := make([]zap.Field, 0, len(context))
	for k, v := range context {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// scrubMessage scrubs the message and fields to remove sensitive data.
func scrubMessage(message string, fields ...zap.Field) {
	// Scrub the log message itself.
	message = scrubber.Scrub(message)

	// Iterate over fields to scrub sensitive values.
	fieldCopy := make([]zap.Field, len(fields))
	copy(fieldCopy, fields)

	for i := range fieldCopy {
		if fieldCopy[i].Key != "" {
			originalValue := fieldCopy[i].String
			switch fieldCopy[i].Type {
			case zapcore.StringType:
				fieldCopy[i] = zap.String(fieldCopy[i].Key, scrubber.Scrub(originalValue))
			case zapcore.ErrorType:
				if err, ok := fieldCopy[i].Interface.(error); ok {
					fieldCopy[i] = zap.String(fieldCopy[i].Key, scrubber.Scrub(err.Error()))
				}
			case zapcore.ReflectType, zapcore.ObjectMarshalerType, zapcore.StringerType:
				fields[i] = zap.String(fields[i].Key, scrubber.Scrub(fmt.Sprintf("%v", fields[i].Interface)))
			default:
				fields[i] = zap.Any(fields[i].Key, scrubber.Scrub(fmt.Sprintf("%v", fields[i].Interface)))
			}
		}
	}
}
