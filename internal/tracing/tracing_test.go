// /observability/tracing/tracing_test.go
package tracing

import (
	"context"
	"errors"
	"testing"

	"github.com/goletan/observability/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// setupTracer configures a test TracerProvider with a SpanRecorder for capturing spans.
func setupTestTracer() *tracetest.SpanRecorder {
	recorder := tracetest.NewSpanRecorder()
	provider := trace.NewTracerProvider(trace.WithSpanProcessor(recorder))
	otel.SetTracerProvider(provider) // Ensure the test provider is set globally.
	tracing.InitTracing(provider)
	return recorder
}

func TestStartSpan(t *testing.T) {
	recorder := setupTestTracer()

	_, span := tracing.StartSpan(context.Background(), "test-span", map[string]interface{}{"attr": "value"})
	defer span.End()

	// Check the recorded spans after ending the span to ensure it's captured.
	span.End()
	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	endedSpan := spans[0]
	if endedSpan.Name() != "test-span" {
		t.Errorf("expected span name 'test-span', got %s", endedSpan.Name())
	}

	// Check if the attribute was set correctly
	expectedAttr := attribute.String("attr", "value")
	if !hasAttribute(endedSpan.Attributes(), expectedAttr) {
		t.Errorf("expected span attribute '%v', got %v", expectedAttr, endedSpan.Attributes())
	}
}

func TestEndSpanWithError(t *testing.T) {
	recorder := setupTestTracer()

	_, span := tracing.StartSpan(context.Background(), "test-span", nil)
	err := errors.New("test error")

	tracing.EndSpan(span, err)

	// Verify that the span was recorded correctly
	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	endedSpan := spans[0]
	// Check if the span has recorded an error status
	if endedSpan.Status().Code != codes.Error {
		t.Errorf("expected span status 'Error', got %v", endedSpan.Status().Code)
	}

	// Check if the error message was recorded in the span
	expectedErrorAttr := attribute.String("error", "true")
	if !hasAttribute(endedSpan.Attributes(), expectedErrorAttr) {
		t.Errorf("expected error attribute 'error: true', got %v", endedSpan.Attributes())
	}
}

func TestEndSpanWithoutError(t *testing.T) {
	recorder := setupTestTracer()

	_, span := tracing.StartSpan(context.Background(), "test-span", nil)
	tracing.EndSpan(span, nil)

	// Verify that the span was recorded correctly
	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	endedSpan := spans[0]
	// Check if the span status was set to OK
	if endedSpan.Status().Code != codes.Ok {
		t.Errorf("expected span status 'Ok', got %v", endedSpan.Status().Code)
	}

	// Check if the success message was recorded in the span
	expectedSuccessAttr := attribute.String("status", "Success")
	if !hasAttribute(endedSpan.Attributes(), expectedSuccessAttr) {
		t.Errorf("expected status attribute 'status: Success', got %v", endedSpan.Attributes())
	}
}

// Helper function to check if a span has a specific attribute
func hasAttribute(attributes []attribute.KeyValue, targetAttr attribute.KeyValue) bool {
	for _, attr := range attributes {
		if attr.Key == targetAttr.Key && attr.Value == targetAttr.Value {
			return true
		}
	}
	return false
}
