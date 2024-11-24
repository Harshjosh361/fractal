package opentele

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace" // Alias for the SDK trace package
	"go.opentelemetry.io/otel/trace"              // OpenTelemetry API trace
)

// InitTracing initializes OpenTelemetry with Jaeger exporter
func InitTracing() (func(), error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Jaeger endpoint (this will be the Jaeger service in the docker-compose network)
	jaegerURL := os.Getenv("JAEGER_URL")
	if jaegerURL == "" {
		jaegerURL = "http://jaeger:14268/api/traces" // Default to the Docker service name 'jaeger'
	}

	// Create a Jaeger exporter to send traces
	traceExporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(jaegerURL),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %v", err)
	}

	// Set up the tracer provider with the exporter and a sampler
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // You can use a custom sampler here
	)

	// Register the TracerProvider globally
	otel.SetTracerProvider(tracerProvider)

	// Return a cleanup function to flush traces when done
	return func() {
		tracerProvider.Shutdown(context.Background())
	}, nil
}

// CreateSpan starts a new span with the provided operation name and adds some attributes
func CreateSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("my-service")
	ctx, span := tracer.Start(ctx, operationName)
	span.SetAttributes(
		// Custom attributes (optional)
		// Example: Set environment as "development"
		attribute.String("env", "development"),
	)
	return ctx, span
}
