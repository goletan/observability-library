// /observability/pkg/observability.go
package observability

import (
	"context"

	"github.com/goletan/observability/internal/config"
	"github.com/goletan/observability/internal/errors"
	"github.com/goletan/observability/internal/logger"
	"github.com/goletan/observability/internal/metrics"
	"github.com/goletan/observability/internal/tracing"
	"github.com/goletan/observability/internal/types"
	"github.com/goletan/observability/internal/utils"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Observability struct {
	Logger  *zap.Logger
	Metrics *metrics.MetricsManager
	Tracer  trace.Tracer
}

func LoadObservabilityConfig(logger *zap.Logger) (*types.ObservabilityConfig, error) {
	return config.LoadObservabilityConfig(logger)
}

func NewObserver(cfg *types.ObservabilityConfig) (*Observability, error) {
	if cfg == nil {
		return nil, errors.WrapError(nil, nil, "observability configuration is nil", 1000, nil)
	}

	log, err := logger.InitLogger(cfg)
	if err != nil {
		return nil, errors.WrapError(nil, err, "Failed to initialize logger", 1001, nil)
	}

	metricsManager, err := metrics.InitMetrics()
	if err != nil {
		return nil, errors.WrapError(log, err, "Failed to initialize metrics", 1002, nil)
	}

	tracer, err := tracing.InitTracing(cfg)
	if err != nil {
		return nil, errors.WrapError(log, err, "Failed to initialize tracing", 1003, nil)
	}

	return &Observability{
		Logger:  log,
		Metrics: metricsManager,
		Tracer:  tracer,
	}, nil
}

func NewScrubber() *utils.Scrubber {
	return utils.NewScrubber()
}

func (o *Observability) Shutdown() error {
	ctx := context.Background()
	if err := tracing.ShutdownTracing(ctx); err != nil {
		o.Logger.Error("Failed to shutdown tracer", zap.Error(err))
		return err
	}
	return nil
}
