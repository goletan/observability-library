// /observability/pkg/observability.go
package observability

import (
	"context"

	"github.com/goletan/observability/config"
	"github.com/goletan/observability/internal/logger"
	"github.com/goletan/observability/internal/metrics"
	"github.com/goletan/observability/internal/tracing"
	"github.com/goletan/observability/internal/utils"
	"github.com/goletan/observability/shared/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Observability struct {
	Logger  *zap.Logger
	Metrics *metrics.MetricsManager
	Tracer  trace.Tracer
}

func NewObserver(cfg *config.ObservabilityConfig) (*Observability, error) {
	log, err := logger.InitLogger(cfg)
	if err != nil {
		return nil, errors.WrapError(nil, err, "Failed to initialize logger", 3001, nil)
	}

	metricsManager, err := metrics.InitMetrics(log)
	if err != nil {
		return nil, errors.WrapError(log, err, "Failed to initialize metrics", 2001, nil)
	}

	tracer, err := tracing.InitTracing(cfg, log)
	if err != nil {
		return nil, errors.WrapError(log, err, "Failed to initialize tracing", 1001, nil)
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
		return errors.WrapError(o.Logger, err, "Failed to shutdown tracer", 1004, nil)
	}
	return nil
}
