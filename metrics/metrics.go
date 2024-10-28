// /observability/metrics/metrics.go
package metrics

import (
	"log"
	"net/http"
	"runtime"
	"time"

	utils "github.com/goletan/observability/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTP Metrics: Track HTTP requests and errors.
var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "api",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
		},
		[]string{"method", "endpoint", "status"},
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
	Timeouts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "timeouts_total",
			Help:      "Counts the number of timeouts in various operations.",
		},
		[]string{"operation", "service"},
	)
	BulkheadLimitReached = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "resilience",
			Name:      "bulkhead_limit_reached_total",
			Help:      "Counts the number of times bulkhead limits have been reached.",
		},
		[]string{"service"},
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

// Security Metrics: Track security events
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

// Messages Metrics: Track production and consumption of messages
var (
	MessagesProduced = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "messages",
			Name:      "goletan_messages_produced_total",
			Help:      "Total number of messages produced.",
		},
		[]string{"event_type", "status", "partition"},
	)
	MessagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "messages",
			Name:      "goletan_messages_consumed_total",
			Help:      "Total number of messages consumed.",
		},
		[]string{"event_type", "status", "partition"},
	)
)

// Application Metrics: Track application related metrics
var (
	AppErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goletan",
			Subsystem: "runtime",
			Name:      "error_count_total",
			Help:      "Total count of errors encountered by the application.",
		},
		[]string{"type", "service", "context"},
	)
	MemoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "goletan",
		Subsystem: "runtime",
		Name:      "memory_usage_bytes",
		Help:      "Current memory usage in bytes.",
	})
	GoroutinesCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "goletan",
		Subsystem: "runtime",
		Name:      "goroutines_count",
		Help:      "Number of goroutines currently running.",
	})
)

// Security Tool: Scrub sensitive data
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
	if err := registerMessagesMetrics(); err != nil {
		log.Printf("Warning: failed to register messages metrics: %v", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Printf("Error starting metrics server: %v", err)
		}
	}()

	go collectRuntimeMetrics()
}

func collectRuntimeMetrics() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		MemoryUsage.Set(float64(m.Alloc))
		GoroutinesCount.Set(float64(runtime.NumGoroutine()))

		time.Sleep(5 * time.Second)
	}
}

func registerHTTPMetrics() error {
	if err := prometheus.Register(RequestDuration); err != nil {
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

	if err := prometheus.Register(Timeouts); err != nil {
		return err
	}

	if err := prometheus.Register(BulkheadLimitReached); err != nil {
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

func registerMessagesMetrics() error {
	if err := prometheus.Register(MessagesProduced); err != nil {
		return err
	}

	if err := prometheus.Register(MessagesConsumed); err != nil {
		return err
	}

	return nil
}

func registerAppMetrics() error {
	if err := prometheus.Register(AppErrorCount); err != nil {
		return err
	}

	if err := prometheus.Register(MemoryUsage); err != nil {
		return err
	}

	if err := prometheus.Register(GoroutinesCount); err != nil {
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
	AppErrorCount.WithLabelValues(scrubbedErrorType, scrubbedService, scrubbedContext).Inc()
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

// IncrementMessagesProduced increments the Kafka produced message counter
func IncrementMessagesProduced(eventType, status, partition string) {
	MessagesProduced.WithLabelValues(eventType, status, partition).Inc()
}

// IncrementMessagesConsumed increments the Kafka consumed message counter
func IncrementMessagesConsumed(eventType, status, partition string) {
	MessagesConsumed.WithLabelValues(eventType, status, partition).Inc()
}

// CountTimeout logs a timeout event for a specific operation and service.
func CountTimeout(operation, service string) {
	Timeouts.WithLabelValues(operation, service).Inc()
}

// CountBulkheadLimitReached logs when a bulkhead limit is reached for a specific service.
func CountBulkheadLimitReached(service string) {
	BulkheadLimitReached.WithLabelValues(service).Inc()
}
