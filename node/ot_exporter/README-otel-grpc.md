Using the OTLP gRPC exporters

This sample has been updated to use the OpenTelemetry gRPC OTLP exporters which send traces and metrics
to the collector over gRPC (default port 4317).

Environment variables
- OTEL_COLLECTOR_TRACES: gRPC URL for traces, e.g. grpc://localhost:4317
- OTEL_COLLECTOR_METRICS: gRPC URL for metrics, e.g. grpc://localhost:4317
- OTEL_SERVICE_NAME: service name

If you don't set the above, the code defaults to grpc://0.0.0.0:4317 which is suitable when
the app runs inside a container and sends data to a collector container.

Installers
Ensure the following packages are installed in the project (they may not be present by default):
- @opentelemetry/exporter-trace-otlp-grpc
- @opentelemetry/exporter-metrics-otlp-grpc

On Windows (PowerShell) run:

npm install @opentelemetry/exporter-trace-otlp-grpc @opentelemetry/exporter-metrics-otlp-grpc

Run the app

Set the env file or export variables and start the server:

node src/server.js

Resolve SSL/TLS "wrong version number" errors

If you see an error like:

PeriodicExportingMetricReader: metrics export failed (error Error: 14 UNAVAILABLE: No connection established. Last error: Error: B8BD0300:error:0A00010B:SSL routines:ssl3_get_record:wrong version number

This usually means your exporter attempted to speak TLS to a plaintext gRPC port (or the URL used an HTTP/HTTPS endpoint). To avoid it:

- Use grpc:// URLs for plaintext gRPC endpoints (e.g. grpc://collector:4317) and the SDK will use insecure credentials.
- If your collector expects TLS, use grpcs:// or configure the collector accordingly.

Verify that your OpenTelemetry Collector is listening on 4317 and configured to receive OTLP/gRPC.
