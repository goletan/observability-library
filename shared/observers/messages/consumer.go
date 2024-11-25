// /observability/shared/observers/messages/consumer.go
package messages

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// StartConsumerSpan starts a span for Kafka message consumption.
func StartConsumerSpan(ctx context.Context, tracer trace.Tracer, msg kafka.Message) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "kafka_consume_message",
		trace.WithAttributes(
			attribute.String("kafka.topic", msg.Topic),
			attribute.Int64("kafka.partition", int64(msg.Partition)),
			attribute.Int64("kafka.offset", msg.Offset),
			attribute.String("operation", "consume"),
		))

	return ctx, span
}
