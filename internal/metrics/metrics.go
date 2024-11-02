// /observability/internal/metrics/metrics.go
package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitMetrics() (*MetricsManager, error) {
	manager := NewManager()

	// Initialize all registered metrics
	if err := manager.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Start the metrics server
	go func() {
		fmt.Println("Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", promhttp.Handler()); err != nil {
			fmt.Printf("Error starting metrics server: %v\n", err)
		}
	}()

	return manager, nil
}
