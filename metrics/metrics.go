// /observability/metrics/metrics.go
package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitMetrics() {
	manager := NewManager()

	// Initialize all registered metrics
	manager.Init()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Printf("Error starting metrics server: %v", err)
		}
	}()
}
