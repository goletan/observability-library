// /observability/metrics/metrics_test.go
package observability_test

import (
	"strings"
	"sync"
	"testing"

	metrics "github.com/goletan/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsCollection(t *testing.T) {
	// Create a custom registry.
	registry := prometheus.NewRegistry()

	// Register the specific collectors we want to test.
	registry.MustRegister(metrics.RequestDuration)
	registry.MustRegister(metrics.ErrorCount)
	registry.MustRegister(metrics.CircuitBreakerStateChange)
	registry.MustRegister(metrics.RetryAttempts)
	registry.MustRegister(metrics.ServiceExecutionDuration)
	registry.MustRegister(metrics.SecurityEvents)

	// Simulate recording some metrics
	metrics.ObserveRequestDuration("GET", "/user", 200, 0.123)
	metrics.ObserveRequestDuration("POST", "/token", 500, 10.123)
	metrics.IncrementErrorCount("password=1234", "secret=service", "apikey=5416dsac")
	metrics.IncrementErrorCount("test_error", "test_service", "test_context")
	metrics.RecordCircuitBreakerStateChange("test_circuit", "closed", "open")
	metrics.CountRetryAttempt("test_operation", "success")
	metrics.ObserveServiceExecution("test_service", "initialize", 0.456)
	metrics.ObserveServiceExecution("critical_service", "process_payment", 999.999)
	metrics.RecordSecurityEvent("failed_login", "auth_service", "high")

	// Concurrent metric updates to test thread safety
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			metrics.CountRetryAttempt("test_concurrent_operation", "retry")
		}()
	}
	wg.Wait()

	// Adjusted expected output for the metrics
	expectedMetrics := `
		# HELP goletan_kernel_http_request_duration_seconds Duration of HTTP requests in seconds.
		# TYPE goletan_kernel_http_request_duration_seconds histogram
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.005"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.01"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.025"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.05"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.1"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.25"} 1
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="0.5"} 1
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="1"} 1
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="2.5"} 1
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="5"} 1
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="10"} 1
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="GET",status="OK",le="+Inf"} 1
		goletan_kernel_http_request_duration_seconds_sum{endpoint="[REDACTED]",method="GET",status="OK"} 0.123
		goletan_kernel_http_request_duration_seconds_count{endpoint="[REDACTED]",method="GET",status="OK"} 1
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.005"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.01"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.025"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.05"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.1"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.25"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="0.5"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="1"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="2.5"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="5"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="10"} 0
		goletan_kernel_http_request_duration_seconds_bucket{endpoint="[REDACTED]",method="POST",status="Internal Server Error",le="+Inf"} 1
		goletan_kernel_http_request_duration_seconds_sum{endpoint="[REDACTED]",method="POST",status="Internal Server Error"} 10.123
		goletan_kernel_http_request_duration_seconds_count{endpoint="[REDACTED]",method="POST",status="Internal Server Error"} 1
		# HELP goletan_kernel_error_count_total Total count of errors encountered by the application.
		# TYPE goletan_kernel_error_count_total counter
		goletan_kernel_error_count_total{context="[REDACTED]",service="[REDACTED]",type="[REDACTED]"} 1
		goletan_kernel_error_count_total{context="test_context",service="test_service",type="test_error"} 1
		# HELP goletan_resilience_circuit_breaker_state_changes_total Tracks state changes in circuit breakers.
		# TYPE goletan_resilience_circuit_breaker_state_changes_total counter
		goletan_resilience_circuit_breaker_state_changes_total{circuit="test_circuit",from="closed",to="open"} 1
		# HELP goletan_resilience_retry_attempts_total Counts the number of retry attempts for operations.
		# TYPE goletan_resilience_retry_attempts_total counter
		goletan_resilience_retry_attempts_total{operation="test_operation",status="success"} 1
		goletan_resilience_retry_attempts_total{operation="test_concurrent_operation",status="retry"} 10
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
		# HELP goletan_security_events_count Counts security-related events like failed authentications.
		# TYPE goletan_security_events_count counter
		goletan_security_events_count{event_type="failed_login",service="auth_service",severity="high"} 1
	`

	// Collect all registered metrics and compare them with the expected metrics.
	if err := testutil.GatherAndCompare(registry, strings.NewReader(expectedMetrics)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}
