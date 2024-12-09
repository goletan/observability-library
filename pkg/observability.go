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

func NewObserver(cfg *config.ObservabilityConfig) (*Observability, error) {
	zapCfg := zap.NewProductionConfig()
	if cfg.Logger.LogLevel != "" {
		if err := zapCfg.Level.UnmarshalText([]byte(cfg.Logger.LogLevel)); err != nil {
			return nil, errors.WrapError(nil, err, "Invalid log level", 1001, nil)
		}
	}

	log, err := logger.NewLogger(&zapCfg)
	if err != nil {
		return nil, errors.WrapError(nil, err, "Failed to initialize logger", 1002, nil)
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

func (o *Observability) Shutdown() error {
	ctx := context.Background()
	if err := tracing.ShutdownTracing(ctx); err != nil {
		o.Logger.Error("Failed to shutdown tracer", zap.Error(err))
		return err
	}

	return nil
}
