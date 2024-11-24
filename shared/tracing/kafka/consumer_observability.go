// /observability/shared/tracing/kafka/consumer_observability.go
package kafka

import (
	"context"

	"github.com/goletan/observability/shared/logger"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ConsumeMessageWithObservability wraps Kafka message consumption with tracing and logging.
func ConsumeMessageWithObservability(ctx context.Context, msg kafka.Message, tracer trace.Tracer, log *logger.ZapLogger) error {
	// Start a new span for Kafka message consumption
	ctx, span := tracer.Start(ctx, "kafka_consume_message",
		trace.WithAttributes(
			attribute.String("kafka.topic", msg.Topic),
			attribute.Int64("kafka.partition", int64(msg.Partition)),
			attribute.Int64("kafka.offset", msg.Offset),
			attribute.String("operation", "consume"),
		))
	defer span.End()

	// Log message consumption start
	log.WithContext(map[string]interface{}{
		"context":   "Consuming message",
		"topic":     msg.Topic,
		"partition": msg.Partition,
		"offset":    msg.Offset,
	})

	// Simulate message processing here (replace with actual logic)
	err := processMessage(ctx, msg)
	if err != nil {
		// Record error details in tracing and logging
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.WithContext(map[string]interface{}{
			"context": "Failed to process message",
			"topic":   msg.Topic,
			"key":     string(msg.Key),
			"error":   err,
		})
		return err
	}

	// Successfully processed the message
	span.SetStatus(codes.Ok, "Message processed successfully")
	log.WithContext(map[string]interface{}{
		"context":   "Message processed successfully",
		"topic":     msg.Topic,
		"partition": msg.Partition,
		"offset":    msg.Offset,
	})

	return nil
}

// processMessage simulates the message processing logic. Replace this with your custom logic.
func processMessage(ctx context.Context, msg kafka.Message) error {
	// For now, just a placeholder function.
	// Replace this with the actual logic for handling the Kafka message.
	return nil
}
