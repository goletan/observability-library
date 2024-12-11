package observability

import (
	"context"

	"github.com/goletan/observability/internal/metrics"
	"github.com/goletan/observability/internal/tracing"
	"github.com/goletan/observability/shared/config"
	"github.com/goletan/observability/shared/errors"
	"github.com/goletan/observability/shared/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Observability struct {
	Logger  *logger.ZapLogger
	Metrics *metrics.Manager
	Tracer  trace.Tracer
}

func NewObserver(cfg *config.ObservabilityConfig, log *logger.ZapLogger) (*Observability, error) {
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

func (o *Observability) Shutdown() error {
	ctx := context.Background()
	if err := tracing.ShutdownTracing(ctx); err != nil {
		o.Logger.Error("Failed to shutdown tracer", zap.Error(err))
		return err
	}

	return nil
}
