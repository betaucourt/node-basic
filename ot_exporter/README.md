# Node.js Fastify OTEL (improved)

This is a copy of the `obs/app/nodejs` example, improved to emit additional metrics and trace-friendly attributes. It's placed in `tmp/nodejs` for experimentation.

How to run:

1. cd tmp/nodejs
2. npm install
3. Set OTEL collector endpoints via environment variables: OTEL_COLLECTOR_ENDPOINT or OTEL_COLLECTOR_TRACES and OTEL_COLLECTOR_METRICS
4. npm start

The app exposes `/test` and `/error` endpoints and emits custom metrics: `http_request_duration_seconds` and `process_rss_bytes`.
