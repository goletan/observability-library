// /observability/metrics/metrics.go
package metrics

import (
	"fmt"
	"log"
	"net/http"

	utils "github.com/goletan/observability/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTP Metrics: Track HTTP requests and errors.
var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "kernel",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
		},
		[]string{"method", "endpoint", "status"},
	)
	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "kernel",
			Name:      "error_count_total",
			Help:      "Total count of errors encountered by the application.",
		},
		[]string{"type", "service", "context"},
	)
)

// Resilience Metrics: Track resilience patterns like retries and circuit breakers.
var (
	CircuitBreakerStateChange = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "circuit_breaker_state_changes_total",
			Help:      "Tracks state changes in circuit breakers.",
		},
		[]string{"circuit", "from", "to"},
	)
	RetryAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "retry_attempts_total",
			Help:      "Counts the number of retry attempts for operations.",
		},
		[]string{"operation", "status"},
	)
)

// Execution Metrics: Track service execution durations.
var (
	ServiceExecutionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "services",
			Name:      "execution_duration_seconds",
			Help:      "Tracks the duration of service execution.",
		},
		[]string{"service", "operation"},
	)
)

// Security Metrics: Track service execution durations.
var (
	SecurityEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "security",
			Name:      "events_count",
			Help:      "Counts security-related events like failed authentications.",
		},
		[]string{"event_type", "service", "severity"},
	)
)

var (
	scrubber = utils.NewScrubber()
)

// InitMetrics initializes Prometheus metrics and starts the metrics server.
func InitMetrics() {
	if err := registerHTTPMetrics(); err != nil {
		log.Printf("Warning: failed to register HTTP metrics: %v", err)
	}
	if err := registerResilienceMetrics(); err != nil {
		log.Printf("Warning: failed to register resilience metrics: %v", err)
	}
	if err := registerExecutionMetrics(); err != nil {
		log.Printf("Warning: failed to register execution metrics: %v", err)
	}
	if err := registerSecurityMetrics(); err != nil {
		log.Printf("Warning: failed to register security metrics: %v", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			fmt.Errorf("Error starting metrics server: %v", err)
		}
	}()
}

func registerHTTPMetrics() error {
	if err := prometheus.Register(RequestDuration); err != nil {
		return err
	}
	if err := prometheus.Register(ErrorCount); err != nil {
		return err
	}
	return nil
}

func registerResilienceMetrics() error {
	if err := prometheus.Register(CircuitBreakerStateChange); err != nil {
		return err
	}
	if err := prometheus.Register(RetryAttempts); err != nil {
		return err
	}
	return nil
}

func registerExecutionMetrics() error {
	if err := prometheus.Register(ServiceExecutionDuration); err != nil {
		return err
	}
	return nil
}

func registerSecurityMetrics() error {
	if err := prometheus.Register(SecurityEvents); err != nil {
		return err
	}
	return nil
}

// ObserveRequestDuration records the duration of HTTP requests.
func ObserveRequestDuration(method, endpoint string, status int, duration float64) {
	scrubbedMethod := scrubber.Scrub(method)
	scrubbedEndpoint := scrubber.Scrub(endpoint)
	RequestDuration.WithLabelValues(scrubbedMethod, scrubbedEndpoint, http.StatusText(status)).Observe(duration)
}

// ObserveServiceExecution records the execution duration of a service operation.
func ObserveServiceExecution(service, operation string, duration float64) {
	scrubbedService := scrubber.Scrub(service)
	scrubbedOperation := scrubber.Scrub(operation)
	ServiceExecutionDuration.WithLabelValues(scrubbedService, scrubbedOperation).Observe(duration)
}

// IncrementErrorCount increments the error counter based on error type, service, and context.
func IncrementErrorCount(errorType, service, context string) {
	scrubbedErrorType := scrubber.Scrub(errorType)
	scrubbedService := scrubber.Scrub(service)
	scrubbedContext := scrubber.Scrub(context)
	ErrorCount.WithLabelValues(scrubbedErrorType, scrubbedService, scrubbedContext).Inc()
}

// RecordCircuitBreakerStateChange logs state changes in the circuit breaker.
func RecordCircuitBreakerStateChange(circuit, from, to string) {
	CircuitBreakerStateChange.WithLabelValues(circuit, from, to).Inc()
}

// CountRetryAttempt logs retry attempts for operations.
func CountRetryAttempt(operation, status string) {
	RetryAttempts.WithLabelValues(operation, status).Inc()
}

// RecordSecurityEvent records a security-related event.
func RecordSecurityEvent(eventType, service, severity string) {
	SecurityEvents.WithLabelValues(eventType, service, severity).Inc()
}
