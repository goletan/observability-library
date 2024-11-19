// /observability/internal/tracing/tracing.go
package tracing

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/goletan/observability/internal/types"
	"github.com/goletan/observability/internal/utils"
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
	scrubber = utils.NewScrubber()
	once     sync.Once
)

// InitTracing initializes OpenTelemetry tracing with an OTLP gRPC exporter.
// It returns the Tracer instance and an error if any occurs during setup.
func InitTracing(cfg *types.ObservabilityConfig, provider ...*sdktrace.TracerProvider) (trace.Tracer, error) {
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
				err = fmt.Errorf("failed to create the collector exporter: %w", err)
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
		log.Println("Tracing initialized successfully")
	})

	if err != nil {
		return nil, err
	}
	return tp.Tracer("goletan-tracer"), nil
}

// StartSpan creates a new span with provided contextual data.
func StartSpan(ctx context.Context, name string, context map[string]interface{}, tracer trace.Tracer) (context.Context, trace.Span) {
	if tracer == nil {
		tracer = otel.Tracer("goletan-tracer")
	}

	ctx, span := tracer.Start(ctx, name)
	maxAttributes := 10
	count := 0
	for k, v := range context {
		if count >= maxAttributes {
			break
		}
		if str, ok := v.(string); ok {
			v = scrubber.Scrub(str)
		}
		span.SetAttributes(attribute.String(k, fmt.Sprintf("%v", v)))
		count++
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

// getServiceName retrieves the service name from the configuration or defaults to "GoletanService".
func getServiceName(cfg *types.ObservabilityConfig) string {
	if cfg.Tracing.ServiceName != "" {
		return cfg.Tracing.ServiceName
	}

	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		return serviceName
	}

	return "GoletanService"
}
