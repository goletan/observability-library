// /observability/logger/logger.go
package observability

import (
	"fmt"
	"sync"

	scrbr "github.com/goletan/observability/scrubber"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger   *zap.Logger
	once     sync.Once
	scrubber = scrbr.NewScrubber()
)

// InitLogger initializes the default logger. It will panic if the logger fails to initialize.
func InitLogger() {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}
	})
}

// SetLogger allows setting a custom logger, primarily used for testing.
func SetLogger(customLogger *zap.Logger) {
	logger = customLogger
}

// ensureLoggerInitialized ensures that the logger is initialized before use.
func ensureLoggerInitialized() {
	if logger == nil {
		InitLogger()
	}
}

// Info logs an info level message with additional fields.
func Info(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Info(message, fields...)
}

// Error logs an error level message with additional fields.
func Error(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Error(message, fields...)
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

// Fatal logs a fatal error and exits the application.
func Fatal(message string, fields ...zap.Field) {
	ensureLoggerInitialized()
	scrubMessage(message, fields...)
	logger.Fatal(message, fields...)
}

// scrubMessage scrubs the message and fields to remove sensitive data.
func scrubMessage(message string, fields ...zap.Field) {
	// Scrub the log message itself.
	message = scrubber.Scrub(message)

	// Iterate over fields to scrub sensitive values.
	for i := range fields {
		originalValue := fields[i].String
		switch fields[i].Type {
		case zapcore.StringType:
			// Scrub the string value of the field.
			fields[i] = zap.String(fields[i].Key, scrubber.Scrub(originalValue))
		case zapcore.ErrorType:
			// Scrub the error value.
			if err, ok := fields[i].Interface.(error); ok {
				fields[i] = zap.String(fields[i].Key, scrubber.Scrub(err.Error()))
			}
		case zapcore.ReflectType, zapcore.ObjectMarshalerType, zapcore.StringerType:
			// Convert complex types to string and scrub.
			fields[i] = zap.String(fields[i].Key, scrubber.Scrub(fmt.Sprintf("%v", fields[i].Interface)))
		default:
			// Scrub the value of other types if they can be converted to a string.
			fields[i] = zap.Any(fields[i].Key, scrubber.Scrub(fmt.Sprintf("%v", fields[i].Interface)))
		}
	}
}
