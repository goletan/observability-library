// /observability/internal/tracing/tracing.go
package tracing

import (
	"context"
	"fmt"
	security "github.com/goletan/security/pkg"
	"sync"

	observability "github.com/goletan/observability/shared/config"

	"github.com/goletan/observability/shared/errors"
	"github.com/goletan/observability/shared/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	scrubber = security.NewScrubber()
	once     sync.Once
)

// InitTracing initializes OpenTelemetry tracing with an OTLP gRPC exporter.
// It returns the Tracer instance and an error if any occurs during setup.
func InitTracing(cfg *observability.ObservabilityConfig, log *logger.ZapLogger, provider ...*sdktrace.TracerProvider) (trace.Tracer, error) {
	var tp *sdktrace.TracerProvider
	var err error
	once.Do(func() {
		if len(provider) > 0 {
			// Use provided TracerProvider, useful for testing or custom setups
			tp = provider[0]
		} else {
			ctx := context.Background()

			// Set up OTLP gRPC exporter
			var exporter *otlptrace.Exporter
			exporter, err = otlptracegrpc.New(ctx)
			if err != nil {
				errors.WrapError(log, err, "Failed to create the OTLP gRPC exporter", 1002, map[string]interface{}{
					"component": "tracing",
					"endpoint":  "grpc",
				})
				return
			}

			// Configure the TracerProvider with a batch span processor and resource attributes
			tp = sdktrace.NewTracerProvider(
				sdktrace.WithBatcher(exporter),
				sdktrace.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(getServiceName(cfg)),
				)),
			)
		}

		// Register the global tracer provider
		otel.SetTracerProvider(tp)
		log.Info("Tracing initialized successfully")
	})

	if err != nil {
		return nil, err
	}
	return tp.Tracer("goletan-tracer"), nil
}

// StartSpanWithMetadata creates a new span with standard metadata.
func StartSpanWithMetadata(cfg *observability.ObservabilityConfig, ctx context.Context, name string, context map[string]interface{}, tracer trace.Tracer) (context.Context, trace.Span) {
	if tracer == nil {
		tracer = otel.Tracer("goletan-tracer")
	}

	ctx, span := tracer.Start(ctx, name,
		trace.WithAttributes(
			attribute.String("service.name", getServiceName(cfg)), // Consistent service name across spans
			attribute.String("environment", cfg.Environment),      // Assuming we have Environment in cfg
			attribute.String("operation.name", name),
		),
	)

	// Add user-defined attributes, scrub sensitive data where needed
	for k, v := range context {
		if str, ok := v.(string); ok {
			v = scrubber.Scrub(str)
		}
		span.SetAttributes(attribute.String(k, fmt.Sprintf("%v", v)))
	}

	return ctx, span
}

// EndSpan finalizes the span and sets its status based on error presence.
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "Success")
	}
	span.End()
}

// ShutdownTracing gracefully shuts down the TracerProvider, exporting all spans.
func ShutdownTracing(ctx context.Context) error {
	tp := otel.GetTracerProvider()
	if provider, ok := tp.(*sdktrace.TracerProvider); ok {
		return provider.Shutdown(ctx)
	}
	return nil
}

// getServiceName retrieves the service name from the configuration.
func getServiceName(cfg *observability.ObservabilityConfig) string {
	if cfg.Tracing.ServiceName != "" {
		return cfg.Tracing.ServiceName
	}

	return "GoletanService"
}
