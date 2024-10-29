// /observability/tracing/tracing.go
package tracing

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/goletan/config"
	"github.com/goletan/observability/types"
	utils "github.com/goletan/observability/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// tracer is the global tracer used throughout the application.
	tracer trace.Tracer
	// scrubber is used to clean sensitive information from context data.
	scrubber = utils.NewScrubber()
	once     sync.Once
	cfg      *types.ObservabilityConfig
)

// InitTracing initializes OpenTelemetry tracing with OTLP exporter or a custom provider.
// Ensure to call ShutdownTracing during application shutdown to flush all spans.
func InitTracing(provider ...*sdktrace.TracerProvider) {
	once.Do(func() {
		var tp *sdktrace.TracerProvider
		if len(provider) > 0 {
			tp = provider[0] // Use the provided TracerProvider for testing
		} else {
			cfg = &types.ObservabilityConfig{}
			err := config.LoadConfig("Observability", cfg, nil)
			if err != nil {
				log.Printf("Warning: failed to load tracing configuration, using defaults: %v", err)
			}

			ctx := context.Background()

			// Set up OTLP gRPC exporter
			exporter, err := otlptracegrpc.New(ctx)
			if err != nil {
				log.Printf("Failed to create the collector exporter: %v", err)
				os.Exit(1)
			}

			// Create a new TracerProvider with a batch span processor and the OTLP exporter
			tp = sdktrace.NewTracerProvider(
				sdktrace.WithBatcher(exporter),
				sdktrace.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(getServiceName()),
				)),
			)
		}

		// Register the global tracer provider
		otel.SetTracerProvider(tp)

		// Set the global tracer to be used in the application
		tracer = tp.Tracer("goletan-tracer")
		log.Println("Tracing initialized successfully")
	})
}

// getServiceName retrieves the service name from the configuration or environment variables.
func getServiceName() string {
	if cfgCache != nil && cfgCache.Tracing.ServiceName != "" {
		return cfgCache.Tracing.ServiceName
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "Goletan"
	}
	return serviceName
}

// GetTracer returns the global tracer for creating spans in the application.
func GetTracer() trace.Tracer {
	return tracer
}

// StartSpan starts a new span with contextual data.
func StartSpan(ctx context.Context, name string, context map[string]interface{}) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, name)
	maxAttributes := 10
	count := 0
	for k, v := range context {
		if count >= maxAttributes {
			break
		}
		if str, ok := v.(string); ok {
			v = scrubber.Scrub(str) // Scrub the context data before adding it to the span
		}
		span.SetAttributes(attribute.String(k, fmt.Sprintf("%v", v)))
		count++
	}

	return ctx, span
}

// EndSpan ends the current span and records error if present.
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "Success")
	}
	span.End()
}

// ShutdownTracing shuts down the tracer provider gracefully, ensuring all spans are exported.
func ShutdownTracing(ctx context.Context) error {
	tp := otel.GetTracerProvider()
	if provider, ok := tp.(*sdktrace.TracerProvider); ok {
		return provider.Shutdown(ctx)
	}
	return nil
}
