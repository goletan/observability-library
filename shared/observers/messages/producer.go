// /observability/shared/observers/messages/producer.go
package messages

import (
	"context"
	"time"

	"github.com/goletan/observability/shared/logger"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ProduceMessageWithObservability produces a Kafka message with tracing and logging capabilities.
func ProduceMessageWithObservability(ctx context.Context, writer *kafka.Writer, msg kafka.Message, tracer trace.Tracer, log *logger.ZapLogger) error {
	// Start a new span for Kafka message production
	ctx, span := tracer.Start(ctx, "kafka_produce_message",
		trace.WithAttributes(
			attribute.String("kafka.topic", msg.Topic),
			attribute.String("kafka.key", string(msg.Key)),
			attribute.String("operation", "produce"),
		))
	defer span.End()

	// Produce message to Kafka
	err := writer.WriteMessages(ctx, msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.WithContext(map[string]interface{}{
			"context": "Failed to produce message",
			"topic":   msg.Topic,
			"key":     string(msg.Key),
		})
		return err
	}

	// Log successful production
	span.SetStatus(codes.Ok, "Success")
	log.WithContext(map[string]interface{}{
		"context": "Produced message successfully",
		"topic":   msg.Topic,
		"key":     string(msg.Key),
		"time":    time.Now().Format(time.RFC3339),
	})
	return nil
}
