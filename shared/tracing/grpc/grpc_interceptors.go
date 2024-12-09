package grpc

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryClientInterceptor is a unary client interceptor for OpenTelemetry tracing.
func UnaryClientInterceptor(tracer trace.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Start a new span for the gRPC call
		ctx, span := tracer.Start(ctx, method)
		defer span.End()

		// Prepare metadata carrier to hold the trace context
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Inject the current context into metadata
		propagation.TraceContext{}.Inject(ctx, propagation.HeaderCarrier(md))

		// Attach the metadata to the outgoing context
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Call the method, propagating the new context
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "Success")
		}

		return err
	}
}

// UnaryServerInterceptor is a unary server interceptor for OpenTelemetry tracing.
func UnaryServerInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract trace information from incoming metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// Extract the tracing context using OpenTelemetry's propagator
			ctx = propagation.TraceContext{}.Extract(ctx, propagation.HeaderCarrier(md))
		}

		// Start a new span for handling the request
		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		// Handle the request and pass the context downstream
		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "Success")
		}

		return resp, err
	}
}
