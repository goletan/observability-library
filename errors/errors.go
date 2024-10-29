package errors

import (
	"fmt"
	"strings"
	"time"

	"github.com/goletan/observability/logger"
	"go.uber.org/zap"
)

// AppError is a structured error that provides a code, message, and additional context.
type AppError struct {
	Code      int
	Message   string
	Err       error
	Timestamp time.Time
	Context   map[string]interface{}
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Message: %s, Error: %v, Timestamp: %s", e.Code, e.Message, e.Err, e.Timestamp.Format(time.RFC3339))
	}
	return fmt.Sprintf("Code: %d, Message: %s, Timestamp: %s", e.Code, e.Message, e.Timestamp.Format(time.RFC3339))
}

// Log the error details using the observability logger.
func logError(err *AppError) {
	logger.Info("Error Occurred", zap.Int("code", err.Code), zap.String("message", err.Message), zap.Any("context", err.Context), zap.String("timestamp", err.Timestamp.Format(time.RFC3339)))
	if err.Err != nil {
		logger.Error("Wrapped Error", zap.Error(err.Err))
	}
}

// WrapError creates a new AppError with an existing error wrapped inside and logs it.
func WrapError(err error, message string, code int, context map[string]interface{}) *AppError {
	appError := &AppError{
		Code:      code,
		Message:   message,
		Err:       err,
		Timestamp: time.Now(),
		Context:   context,
	}
	logError(appError)
	return appError
}

// NewError creates a new AppError without wrapping an existing error and logs it.
func NewError(message string, code int, context map[string]interface{}) *AppError {
	appError := &AppError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Context:   context,
	}
	logError(appError)
	return appError
}

// IsRetryable checks if the error is considered retryable based on error code or message patterns.
func IsRetryable(err error) bool {
	// Attempt to cast the error to AppError type
	appErr, ok := err.(*AppError)
	if !ok {
		// If it's not an AppError, you might have other transient errors in your application
		return false
	}

	// Define a list of retryable error codes (e.g., 502 Bad Gateway, 503 Service Unavailable)
	retryableCodes := map[int]bool{
		502: true,
		503: true,
		504: true, // Gateway Timeout
		408: true, // Request Timeout
	}

	// Check if the error code is in the retryable codes map
	if retryableCodes[appErr.Code] {
		return true
	}

	// Optionally, check for specific keywords in the error message indicating a transient error
	retryableMessages := []string{
		"timeout",       // General timeouts
		"temporarily",   // Temporarily unavailable
		"connection",    // Connection issues
		"rate limited",  // Rate limiting errors
		"unavailable",   // Service unavailable
		"network error", // Network related errors
	}

	for _, msg := range retryableMessages {
		if containsIgnoreCase(appErr.Message, msg) {
			return true
		}
	}

	return false
}

// containsIgnoreCase checks if a substring exists in a string regardless of case
func containsIgnoreCase(source, substring string) bool {
	return strings.Contains(strings.ToLower(source), strings.ToLower(substring))
}
