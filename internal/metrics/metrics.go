package metrics

import (
	logger "github.com/goletan/logger-library/pkg"
	"github.com/goletan/observability-library/internal/types"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"

	"github.com/goletan/observability-library/shared/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitMetrics(cfg *types.ObservabilityConfig, log *logger.ZapLogger) (*Manager, error) {
	manager := NewManager()

	// Initialize all registered metrics
	if err := manager.Init(log); err != nil {
		return nil, errors.WrapError(log, err, "failed to initialize metrics", 2001, map[string]interface{}{
			"component": "metrics",
			"endpoint":  "init",
		})
	}

	metricsAddress := cfg.Metrics.Address
	if metricsAddress == "" || !strings.HasPrefix(metricsAddress, ":") {
		log.Warn("No metrics address or port provided, using default (:2119)")
		metricsAddress = ":2119"
	}

	// Start the metrics server
	go func() {
		log.Info("Starting metrics server on", zap.String("address", metricsAddress))
		if err := http.ListenAndServe(metricsAddress, promhttp.Handler()); err != nil {
			wrappedErr := errors.WrapError(log, err, "failed to start metrics server", 2002, map[string]interface{}{
				"component": "metrics",
				"endpoint":  "ListenAndServe",
			})
			log.Error("Critical: Metrics server failed to start", zap.Error(wrappedErr))
			os.Exit(1)
		}
	}()

	return manager, nil
}
