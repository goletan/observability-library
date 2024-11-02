// /observability/pkg/observability.go
package pkg

import (
	"context"

	"github.com/goletan/observability/internal/errors"
	"github.com/goletan/observability/internal/logger"
	"github.com/goletan/observability/internal/metrics"
	"github.com/goletan/observability/internal/tracing"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Observability struct {
	Logger  *zap.Logger
	Metrics *metrics.MetricsManager
	Tracer  trace.Tracer
}

func NewObserver() (*Observability, error) {
	// Initialize each module independently
	log, err := logger.InitLogger()
	if err != nil {
		return nil, errors.WrapError(nil, err, "Failed to initialize logger", 1001, nil)
	}

	metricsManager, err := metrics.InitMetrics()
	if err != nil {
		return nil, errors.WrapError(log, err, "Failed to initialize metrics", 1002, nil)
	}

	tracer, err := tracing.InitTracing()
	if err != nil {
		return nil, errors.WrapError(log, err, "Failed to initialize tracing", 1003, nil)
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
