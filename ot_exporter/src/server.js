import 'dotenv/config';
import Fastify from 'fastify';
import { sdk, meter, httpRequestDuration, recordHttpDuration } from './otel.js';
import { trace, context, SpanStatusCode } from '@opentelemetry/api';

const fastify = Fastify({
  logger: {
    level: process.env.LOG_LEVEL || 'info',
    mixin() {
      const span = trace.getSpan(context.active());
      const traceId = span?.spanContext()?.traceId;
      return { trace_id: traceId };
    }
  }
});

fastify.get('/test', async (request, reply) => {
  const span = trace.getSpan(context.active());
  const start = Date.now();
  try {
    // Simulate some work
    const result = { kirikou: 'est petit' };
    const durationMs = Date.now() - start;

    // record histogram in seconds with attributes; pass current context so the Meter may
    // attach an exemplar linked to the active trace/span.
    recordHttpDuration(durationMs / 1000, {
      'http.method': request.method,
      'http.route': request.routerPath || request.raw.url,
      'http.status_code': '200'
    }, context.active());

    span?.addEvent('request.handled', {
      'http.method': request.method,
      'http.route': request.routerPath || request.raw.url,
      'http.status_code': 200,
      'http.duration_ms': durationMs,
      'app.response_body': JSON.stringify(result),
    });
    return result;
  } catch (err) {
    const durationMs = Date.now() - start;
    span?.setStatus({ code: SpanStatusCode.ERROR, message: String(err?.message || err) });
    span?.recordException(err);
    recordHttpDuration(durationMs / 1000, {
      'http.method': request.method,
      'http.route': request.routerPath || request.raw.url,
      'http.status_code': '500'
    }, context.active());
    span?.addEvent('request.error', {
      'http.method': request.method,
      'http.route': request.routerPath || request.raw.url,
      'http.status_code': 500,
      'http.duration_ms': durationMs,
      'error.message': String(err?.message || err),
    });
    throw err;
  }
});

fastify.get('/error', async (request, reply) => {
  const span = trace.getSpan(context.active());
  span?.addEvent('controller.enter', {
    'http.method': request.method,
    'http.route': request.routerPath || request.raw.url,
  });

  try {
    throw new Error('simulated controller failure');
  } catch (err) {
    span?.setStatus({ code: SpanStatusCode.ERROR, message: String(err?.message || err) });
    span?.recordException(err);
    span?.addEvent('controller.error', {
      'error.message': String(err?.message || err),
      'http.status_code': 500,
    });
    reply.code(500);
    return { error: 'internal_server_error', message: String(err?.message || err) };
  }
});

const port = process.env.PORT || 8080;

process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Rejection at:', promise, 'reason:', reason);
});
process.on('uncaughtException', (err) => {
  console.error('Uncaught Exception:', err);
  process.exit(1);
});

async function start() {
  try {
    await fastify.listen({ port, host: '0.0.0.0' });
    console.info(`Server listening on ${port}`);
  } catch (err) {
    console.error('Error starting server:', err);
    process.exit(1);
  }
}

start();
