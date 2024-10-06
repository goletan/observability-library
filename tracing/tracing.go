// /observability/tracing/tracing.go
package observability

import (
	"context"
	"fmt"
	"log"

	scrbr "github.com/goletan/observability/scrubber"
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
	tracer   trace.Tracer
	scrubber = scrbr.NewScrubber()
)

// InitTracing initializes OpenTelemetry tracing with OTLP exporter or a custom provider.
func InitTracing(provider ...*sdktrace.TracerProvider) {
	var tp *sdktrace.TracerProvider
	if len(provider) > 0 {
		tp = provider[0] // Use the provided TracerProvider for testing
	} else {
		ctx := context.Background()

		// Set up OTLP gRPC exporter
		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			fmt.Errorf("Failed to create the collector exporter: %v", err)
		}

		// Create a new TracerProvider with a batch span processor and the OTLP exporter
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("Goletan"),
			)),
		)
	}

	// Register the global tracer provider
	otel.SetTracerProvider(tp)

	// Set the global tracer to be used in the application
	tracer = tp.Tracer("goletan-tracer")
	log.Println("Tracing initialized successfully")
}

// GetTracer returns the global tracer for creating spans in the application.
func GetTracer() trace.Tracer {
	return tracer
}

// StartSpan starts a new span with contextual data.
func StartSpan(ctx context.Context, name string, context map[string]interface{}) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, name)
	// Convert the context map to span attributes
	for k, v := range context {
		if str, ok := v.(string); ok {
			v = scrubber.Scrub(str) // Scrub the context data before adding it to the span
		}
		span.SetAttributes(attribute.String(k, fmt.Sprintf("%v", v)))
	}
	return ctx, span
}

// EndSpan ends the current span and records error if present.
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "true")) // Mark the error explicitly
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetAttributes(attribute.String("status", "Success")) // Mark success explicitly
		span.SetStatus(codes.Ok, "Success")
	}
	span.End()
}
