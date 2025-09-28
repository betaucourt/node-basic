// improved, simplified otel.js
import { NodeSDK } from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import FastifyOtelInstrumentation from '@fastify/otel';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-grpc';
import metricsPkg from '@opentelemetry/sdk-metrics';
import { trace as apiTrace, metrics, context as otelContext } from '@opentelemetry/api';
import { credentials as grpcCredentials } from '@grpc/grpc-js';
import { ExpressInstrumentation } from '@opentelemetry/instrumentation-express';
import { FsInstrumentation } from '@opentelemetry/instrumentation-fs';
import * as resourcesPkg from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';

const { PeriodicExportingMetricReader, MeterProvider, TraceIdRatioExemplarSampler } = metricsPkg;

// Create a TraceId-ratio exemplar sampler: use SDK's implementation when present, else a small fallback
function createTraceIdRatioExemplarSampler(ratio) {
  if (typeof TraceIdRatioExemplarSampler === 'function') {
    try {
      return new TraceIdRatioExemplarSampler(ratio);
    } catch (e) {
      // fallthrough to fallback
    }
  }

  // Simple probabilistic fallback that attaches trace/span ids when present in the context
  return new (class {
    constructor(_ratio) {
      this._ratio = Math.max(0, Math.min(1, Number(_ratio) || 0));
    }
    shouldSample(_value, _attributes, _context) {
      const sampled = Math.random() < this._ratio;
      if (!sampled) return { sampled: false, shouldSample: false };

      let traceId, spanId;
      try {
        const span = apiTrace.getSpan(_context);
        if (span && typeof span.spanContext === 'function') {
          const sc = span.spanContext();
          traceId = sc?.traceId;
          spanId = sc?.spanId;
        }
      } catch (e) {
        // ignore
      }

      return { sampled: true, shouldSample: true, traceId, spanId, filteredAttributes: Object.assign({}, _attributes || {}) };
    }
  })(ratio);
}

const defaultGrpcEndpoint = 'grpc://0.0.0.0:4317';
const traceRaw = defaultGrpcEndpoint;
const metricsRaw = defaultGrpcEndpoint;

const traceExporterOpts = { url: traceRaw };
const metricExporterOpts = { url: metricsRaw };
try { if (traceRaw && traceRaw.startsWith('grpc://')) traceExporterOpts.credentials = grpcCredentials.createInsecure(); } catch (e) {}
try { if (metricsRaw && metricsRaw.startsWith('grpc://')) metricExporterOpts.credentials = grpcCredentials.createInsecure(); } catch (e) {}

const traceExporter = new OTLPTraceExporter(traceExporterOpts);
const metricExporter = new OTLPMetricExporter(metricExporterOpts);

const serviceName = process.env.OTEL_SERVICE_NAME || 'nodejs-service-improved';
const serviceVersion = process.env.SERVICE_VERSION || process.env.npm_package_version || '0.1.0';

const ResourceCandidate = resourcesPkg.Resource || resourcesPkg.default?.Resource || resourcesPkg.default || resourcesPkg;
let resource;
const resourceAttrs = {
  [SemanticResourceAttributes.SERVICE_NAME]: serviceName,
  [SemanticResourceAttributes.SERVICE_VERSION]: serviceVersion,
  'deployment.environment': process.env.DEPLOYMENT_ENV || 'local',
};
if (typeof ResourceCandidate === 'function') {
  try { resource = new ResourceCandidate(resourceAttrs); } catch (e) { if (typeof ResourceCandidate.create === 'function') resource = ResourceCandidate.create(resourceAttrs); }
} else if (ResourceCandidate && typeof ResourceCandidate.create === 'function') {
  resource = ResourceCandidate.create(resourceAttrs);
}

const metricReader = new PeriodicExportingMetricReader({ exporter: metricExporter, exportIntervalMillis: 10000 });
const exemplarRatio = Number(process.env.OTEL_EXEMPLAR_SAMPLE_RATIO ?? process.env.EXEMPLAR_SAMPLE_RATIO ?? 0.1);

const meterProvider = new MeterProvider({ resource, exemplarSampler: createTraceIdRatioExemplarSampler(exemplarRatio) });
// Register metric reader with best-effort compatibility
try {
  if (typeof meterProvider.addMetricReader === 'function') meterProvider.addMetricReader(metricReader);
  else if (typeof meterProvider.registerMetricReader === 'function') meterProvider.registerMetricReader(metricReader);
  else if (typeof meterProvider.addReader === 'function') meterProvider.addReader(metricReader);
} catch (e) { console.log('Failed to register metric reader', e); }

const instrumentations = [ getNodeAutoInstrumentations(), new FastifyOtelInstrumentation({ registerOnInitialization: true }), new ExpressInstrumentation(), new FsInstrumentation() ];

const sdk = new NodeSDK({ resource, traceExporter, instrumentations, meterProvider, metricReaders: [metricReader] });
console.log('Starting improved OpenTelemetry SDK');
sdk.start();

const meter = metrics.getMeter('nodejs-improvements', serviceVersion);
export const httpRequestDuration = meter.createHistogram('http_request_duration_seconds', { description: 'Duration of HTTP requests', unit: 's' });

export function recordHttpDuration(valueSeconds, attributes = {}, ctx = undefined) {
  if (ctx) return otelContext.with(ctx, () => httpRequestDuration.record(valueSeconds, attributes));
  return httpRequestDuration.record(valueSeconds, attributes);
}

meter.createObservableGauge('process_rss_bytes', { description: 'Process resident set size', unit: 'By' }).addCallback((observableResult) => {
  try { const rss = process?.memoryUsage?.().rss || 0; observableResult.observe(rss, { 'process.pid': String(process.pid) }); } catch (e) { }
});

process.on('SIGTERM', async () => { try { await sdk.shutdown(); } catch (e) {} });

export { sdk, meter };
