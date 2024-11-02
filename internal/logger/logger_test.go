// /observability/logger/logger_test.go
package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestInfoLog(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	logger.SetLogger(logger)

	logger.Info("Info log test", zap.String("key", "value"))

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}

	logEntry := logs.All()[0]
	if logEntry.Message != "Info log test" {
		t.Errorf("unexpected log message: %s", logEntry.Message)
	}
	if logEntry.Level != zap.InfoLevel {
		t.Errorf("expected level Info, got %s", logEntry.Level)
	}
	if logEntry.Context[0].Key != "key" || logEntry.Context[0].String != "value" {
		t.Errorf("unexpected log field: %v", logEntry.Context[0])
	}
}

func TestErrorLog(t *testing.T) {
	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	logger.SetLogger(logger)

	logger.Error("Error log test", zap.String("error", "something went wrong"))

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}

	logEntry := logs.All()[0]
	if logEntry.Message != "Error log test" {
		t.Errorf("unexpected log message: %s", logEntry.Message)
	}
	if logEntry.Level != zap.ErrorLevel {
		t.Errorf("expected level Error, got %s", logEntry.Level)
	}
	if logEntry.Context[0].Key != "error" || logEntry.Context[0].String != "something went wrong" {
		t.Errorf("unexpected log field: %v", logEntry.Context[0])
	}
}

func TestScrubbedLog(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	logger.SetLogger(logger)

	// Log a message with sensitive data
	logger.Info("Sensitive info", zap.String("password", "mysecretpassword"))

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}

	logEntry := logs.All()[0]
	// Check if the message and context fields are scrubbed correctly
	if logEntry.Message != "Sensitive info" {
		t.Errorf("unexpected log message: %s", logEntry.Message)
	}
	if logEntry.Context[0].Key != "password" || logEntry.Context[0].String != "[REDACTED]" {
		t.Errorf("sensitive data not scrubbed: %v", logEntry.Context[0])
	}
}
