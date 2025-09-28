package middleware

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	appotel "music-app/src/otel"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// HTTPTelemetryMiddleware wraps handlers with OpenTelemetry instrumentation
func HTTPTelemetryMiddleware(next http.HandlerFunc, operationName string) http.HandlerFunc {
	// Wrap with automatic HTTP tracing
	instrumentedHandler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom span for the operation
		tracer := otel.Tracer("music-app")
		ctx, span := tracer.Start(r.Context(), operationName)
		defer span.End()

		// Add request attributes to span
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.route", r.URL.Path),
		)

		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))

		// Record metrics after request completion
		duration := time.Since(start).Seconds()

		// Record HTTP request duration
		if appotel.HTTPRequestDuration != nil {
			appotel.HTTPRequestDuration.Record(r.Context(), duration,
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("route", r.URL.Path),
				))
		}

		// Record process RSS memory periodically (simple approach)
		recordMemoryMetrics(r.Context())
	}), operationName)

	return instrumentedHandler.ServeHTTP
}

// recordMemoryMetrics records current memory usage
func recordMemoryMetrics(ctx context.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if appotel.ProcessRSSBytes != nil {
		// Record RSS memory in bytes
		appotel.ProcessRSSBytes.Record(ctx, int64(m.Sys),
			metric.WithAttributes(
				attribute.String("process.pid", strconv.Itoa(os.Getpid())),
			))
	}
}
