package metrics

import (
	"github.com/goletan/observability/internal/logger"
	"net/http"

	"github.com/goletan/observability/shared/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitMetrics(log *logger.ZapLogger) (*Manager, error) {
	manager := NewManager()

	// Initialize all registered metrics
	if err := manager.Init(log); err != nil {
		return nil, errors.WrapError(log, err, "failed to initialize metrics", 2001, map[string]interface{}{
			"component": "metrics",
			"endpoint":  "init",
		})
	}

	// Start the metrics server
	go func() {
		log.Info("Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", promhttp.Handler()); err != nil {
			err = errors.WrapError(log, err, "failed to start metrics server", 2002, map[string]interface{}{
				"component": "metrics",
				"endpoint":  "ListenAndServe",
			})
		}
	}()

	return manager, nil
}
