// /observability/internal/metrics/metrics_test.go
package metrics

import (
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSecurityEventMetrics(t *testing.T) {
	// Create a custom securityRegistry.
	securityRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	securityRegistry.MustRegister(SecurityEvents)

	// Simulate recording some metrics
	RecordSecurityEvent("failed_login", "auth_service", "high")

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_security_events_count Counts security-related events like failed authentications.
		# TYPE goletan_security_events_count counter
		goletan_security_events_count{event_type="failed_login",service="auth_service",severity="high"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(securityRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestRequestDurationMetrics(t *testing.T) {
	// Create a custom requestRegistry.
	requestRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	requestRegistry.MustRegister(RequestDuration)

	// Simulate recording some metrics
	ObserveRequestDuration("GET", "/user", 200, 0.123)
	ObserveRequestDuration("POST", "/token", 500, 10.123)

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_api_http_request_duration_seconds Duration of HTTP requests in seconds.
		# TYPE goletan_api_http_request_duration_seconds histogram
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.005"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.01"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.025"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.05"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.1"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.25"} 1
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.5"} 1
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="1"} 1
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="2.5"} 1
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="5"} 1
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="10"} 1
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="+Inf"} 1
		goletan_api_http_request_duration_seconds_sum{endpoint="[REDACTED]",method="GET",status="OK"} 0.123
		goletan_api_http_request_duration_seconds_count{endpoint="[REDACTED]",method="GET",status="OK"} 1
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.005"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.01"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.025"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.05"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.1"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.25"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.5"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="1"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="2.5"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="5"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="10"} 0
		goletan_api_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="+Inf"} 1
		goletan_api_http_request_duration_seconds_sum{endpoint="[REDACTED]",method="POST",status="Internal Server Error"} 10.123
		goletan_api_http_request_duration_seconds_count{endpoint="[REDACTED]",method="POST",status="Internal Server Error"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(requestRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestRetryAttemptsMetrics(t *testing.T) {
	// Create a custom retryRegistry.
	retryRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	retryRegistry.MustRegister(RetryAttempts)

	// Simulate recording some metrics
	CountRetryAttempt("test_operation", "success")

	// Concurrent metric updates to test thread safety
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			CountRetryAttempt("test_concurrent_operation", "retry")
		}()
	}
	wg.Wait()

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_resilience_retry_attempts_total Counts the number of retry attempts for operations.
		# TYPE goletan_resilience_retry_attempts_total counter
		goletan_resilience_retry_attempts_total{operation="test_operation",status="success"} 1
		goletan_resilience_retry_attempts_total{operation="test_concurrent_operation",status="retry"} 10
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(retryRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestCircuitBreakerStateChangeMetrics(t *testing.T) {
	// Create a custom circuitBreakerRegistry.
	circuitBreakerRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	circuitBreakerRegistry.MustRegister(CircuitBreakerStateChange)

	// Simulate recording some metrics
	RecordCircuitBreakerStateChange("test_circuit", "closed", "open")

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_resilience_circuit_breaker_state_changes_total Tracks state changes in circuit breakers.
		# TYPE goletan_resilience_circuit_breaker_state_changes_total counter
		goletan_resilience_circuit_breaker_state_changes_total{circuit="test_circuit",from="closed",to="open"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(circuitBreakerRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestServiceExecutionDurationMetrics(t *testing.T) {
	// Create a custom serviceRegistry.
	serviceRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	serviceRegistry.MustRegister(ServiceExecutionDuration)

	// Simulate recording some metrics
	ObserveServiceExecution("test_service", "initialize", 0.456)
	ObserveServiceExecution("critical_service", "process_payment", 999.999)

	// Adjusted expected output for the metrics
	expectedMetrics := `
	# HELP goletan_services_execution_duration_seconds Tracks the duration of service execution.
         # TYPE goletan_services_execution_duration_seconds histogram
         goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="0.005"} 0
         goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="0.01"} 0
         goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="0.025"} 0
         goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="0.05"} 0
         goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="0.1"} 0
         goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="0.25"} 0
        goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="0.5"} 1
        goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="1"} 1
        goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="2.5"} 1
        goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="5"} 1
        goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="10"} 1
        goletan_services_execution_duration_seconds_bucket{operation="initialize",service="test_service",le="+Inf"} 1
        goletan_services_execution_duration_seconds_sum{operation="initialize",service="test_service"} 0.456
        goletan_services_execution_duration_seconds_count{operation="initialize",service="test_service"} 1
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="0.005"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="0.01"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="0.025"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="0.05"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="0.1"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="0.25"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="0.5"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="1"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="2.5"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="5"} 0
         goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="10"} 0
        goletan_services_execution_duration_seconds_bucket{operation="process_payment",service="critical_service",le="+Inf"} 1
        goletan_services_execution_duration_seconds_sum{operation="process_payment",service="critical_service"} 999.999
        goletan_services_execution_duration_seconds_count{operation="process_payment",service="critical_service"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(serviceRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestTimeoutsMetrics(t *testing.T) {
	// Create a custom timeoutsRegistry.
	timeoutsRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	timeoutsRegistry.MustRegister(Timeouts)

	// Simulate recording some metrics
	CountTimeout("test_operation", "test_service")

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_resilience_timeouts_total Counts the number of timeouts in various operations.
        # TYPE goletan_resilience_timeouts_total counter
        goletan_resilience_timeouts_total{operation="test_operation",service="test_service"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(timeoutsRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestBulkheadLimitReachedMetrics(t *testing.T) {
	// Create a custom bulkheadRegistry.
	bulkheadRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	bulkheadRegistry.MustRegister(BulkheadLimitReached)

	// Simulate recording some metrics
	CountBulkheadLimitReached("auth")

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_resilience_bulkhead_limit_reached_total Counts the number of times bulkhead limits have been reached.
        # TYPE goletan_resilience_bulkhead_limit_reached_total counter
        goletan_resilience_bulkhead_limit_reached_total{service="auth"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(bulkheadRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestAppErrorMetrics(t *testing.T) {
	// Create a custom appRegistry.
	appRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	appRegistry.MustRegister(AppErrorCount)

	// Simulate recording some metrics
	IncrementErrorCount("password=1234", "secret=service", "apikey=5416dsac")
	IncrementErrorCount("test_error", "test_service", "test_context")

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_runtime_error_count_total Total count of errors encountered by the application.
		# TYPE goletan_runtime_error_count_total counter
		goletan_runtime_error_count_total{context="[REDACTED]",service="[REDACTED]",type="[REDACTED]"} 1
		goletan_runtime_error_count_total{context="test_context",service="test_service",type="test_error"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(appRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestMessagesMetrics(t *testing.T) {
	// Create a custom messagesRegistry.
	messagesRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	messagesRegistry.MustRegister(MessagesProduced)
	messagesRegistry.MustRegister(MessagesConsumed)
	messagesRegistry.MustRegister(GoroutinesCount)

	// Simulate recording some metrics
	IncrementMessagesProduced("user_created", "OK", "partition")
	IncrementMessagesConsumed("user_created", "OK", "partition")

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_messages_produced_total Total number of messages produced.
		# TYPE goletan_messages_produced_total counter
		goletan_messages_produced_total{event_type="user_created",partition="partition",status="OK"} 1
		# HELP goletan_messages_consumed_total Total number of messages consumed.
		# TYPE goletan_messages_consumed_total counter
		goletan_messages_consumed_total{event_type="user_created",partition="partition",status="OK"} 1
		# HELP goletan_runtime_goroutines_count Number of goroutines currently running.
    	# TYPE goletan_runtime_goroutines_count gauge
        goletan_runtime_goroutines_count 0
    `

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(messagesRegistry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestGoroutineMetrics(t *testing.T) {
	// Create a custom goroutineRegistry.
	goroutineRegistry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	goroutineRegistry.MustRegister(MemoryUsage)
	goroutineRegistry.MustRegister(GoroutinesCount)

	// Create a channel to stop the metric collection after the test.
	done := make(chan bool)
	defer close(done)

	// Start collecting runtime metrics.
	go collectRuntimeMetrics(done)

	// Wait for a few cycles of metric collection.
	time.Sleep(6 * time.Second)

	// Get the actual number of goroutines at this point.
	expectedGoroutines := float64(runtime.NumGoroutine())

	// Allow some flexibility, e.g., a 20% difference.
	tolerance := 0.2
	minGoroutines := expectedGoroutines * (1 - tolerance)
	maxGoroutines := expectedGoroutines * (1 + tolerance)

	// Collect and compare goroutine count.
	collectedGoroutines := testutil.ToFloat64(GoroutinesCount)
	if collectedGoroutines < minGoroutines || collectedGoroutines > maxGoroutines {
		t.Errorf("unexpected goroutines count: got %v, want within range [%v, %v]", collectedGoroutines, minGoroutines, maxGoroutines)
	}

	// Get the actual memory usage.
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	expectedMemoryUsage := float64(m.Alloc)

	// Allow a similar tolerance for memory usage.
	minMemoryUsage := expectedMemoryUsage * (1 - tolerance)
	maxMemoryUsage := expectedMemoryUsage * (1 + tolerance)

	// Collect and compare memory usage.
	collectedMemoryUsage := testutil.ToFloat64(MemoryUsage)
	if collectedMemoryUsage < minMemoryUsage || collectedMemoryUsage > maxMemoryUsage {
		t.Errorf("unexpected memory usage: got %v, want within range [%v, %v]", collectedMemoryUsage, minMemoryUsage, maxMemoryUsage)
	}

	// Signal the goroutine to stop.
	done <- true
}
