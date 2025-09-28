package otel

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// Meter for custom metrics
	Meter metric.Meter

	// Custom metrics
	HTTPRequestDuration metric.Float64Histogram
	DBQueryDuration     metric.Float64Histogram
	ProcessRSSBytes     metric.Int64Gauge
)

// InitOTel initializes OpenTelemetry with OTLP exporters
func InitOTel(serviceName, serviceVersion string) func(context.Context) error {
	ctx := context.Background()

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// OTLP gRPC endpoint (same as Node.js app)
	endpoint := "0.0.0.0:4317"

	// Set up trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	// Set up trace provider
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set up metric exporter
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlpmetricgrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		log.Fatalf("failed to create metric exporter: %v", err)
	}

	// Set up metric provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
			sdkmetric.WithInterval(10*time.Second))),
	)
	otel.SetMeterProvider(meterProvider)

	// Set up propagators
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Create meter for custom metrics
	Meter = otel.Meter("music-app", metric.WithInstrumentationVersion(serviceVersion))

	// Initialize custom metrics
	initMetrics()

	fmt.Println("OpenTelemetry initialized successfully")

	// Return shutdown function
	return func(ctx context.Context) error {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
		return nil
	}
}

// initMetrics creates custom metrics instruments
func initMetrics() {
	var err error

	// HTTP request duration histogram
	HTTPRequestDuration, err = Meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Printf("failed to create http_request_duration metric: %v", err)
	}

	// Database query duration histogram
	DBQueryDuration, err = Meter.Float64Histogram(
		"db_query_duration_seconds",
		metric.WithDescription("Duration of database queries"),
		metric.WithUnit("s"),
	)
	if err != nil {
		log.Printf("failed to create db_query_duration metric: %v", err)
	}

	// Process RSS memory gauge
	ProcessRSSBytes, err = Meter.Int64Gauge(
		"process_rss_bytes",
		metric.WithDescription("Process resident set size"),
		metric.WithUnit("By"),
	)
	if err != nil {
		log.Printf("failed to create process_rss_bytes metric: %v", err)
	}
}
