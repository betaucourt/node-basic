package database

import (
	"context"
	"database/sql"
	"time"

	appotel "music-app/src/otel"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedDB wraps sql.DB with OpenTelemetry instrumentation
type InstrumentedDB struct {
	*sql.DB
	tracer trace.Tracer
}

// NewInstrumentedDB creates a new instrumented database connection
func NewInstrumentedDB(db *sql.DB) *InstrumentedDB {
	return &InstrumentedDB{
		DB:     db,
		tracer: otel.Tracer("music-app-db"),
	}
}

// Query executes a query with instrumentation
func (idb *InstrumentedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	// Create span for database query
	ctx, span := idb.tracer.Start(ctx, "db.query")
	defer span.End()

	// Add attributes to span
	span.SetAttributes(
		attribute.String("db.operation", "query"),
		attribute.String("db.sql", query),
	)

	// Execute the query
	rows, err := idb.DB.QueryContext(ctx, query, args...)

	// Record metrics
	duration := time.Since(start).Seconds()
	if appotel.DBQueryDuration != nil {
		appotel.DBQueryDuration.Record(ctx, duration,
			metric.WithAttributes(
				attribute.String("operation", "query"),
				attribute.Bool("error", err != nil),
			))
	}

	if err != nil {
		span.RecordError(err)
	}

	return rows, err
}

// QueryRow executes a single-row query with instrumentation
func (idb *InstrumentedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()

	// Create span for database query
	ctx, span := idb.tracer.Start(ctx, "db.query_row")
	defer span.End()

	// Add attributes to span
	span.SetAttributes(
		attribute.String("db.operation", "query_row"),
		attribute.String("db.sql", query),
	)

	// Execute the query
	row := idb.DB.QueryRowContext(ctx, query, args...)

	// Record metrics
	duration := time.Since(start).Seconds()
	if appotel.DBQueryDuration != nil {
		appotel.DBQueryDuration.Record(ctx, duration,
			metric.WithAttributes(
				attribute.String("operation", "query_row"),
				attribute.Bool("error", false), // QueryRow doesn't return error directly
			))
	}

	return row
}

// Exec executes a statement with instrumentation
func (idb *InstrumentedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()

	// Create span for database exec
	ctx, span := idb.tracer.Start(ctx, "db.exec")
	defer span.End()

	// Add attributes to span
	span.SetAttributes(
		attribute.String("db.operation", "exec"),
		attribute.String("db.sql", query),
	)

	// Execute the statement
	result, err := idb.DB.ExecContext(ctx, query, args...)

	// Record metrics
	duration := time.Since(start).Seconds()
	if appotel.DBQueryDuration != nil {
		appotel.DBQueryDuration.Record(ctx, duration,
			metric.WithAttributes(
				attribute.String("operation", "exec"),
				attribute.Bool("error", err != nil),
			))
	}

	if err != nil {
		span.RecordError(err)
	}

	return result, err
}
