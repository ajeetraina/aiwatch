# AIWatch Observability

This directory contains documentation and configuration for the observability features added to the AIWatch platform.

## Observability Features

The following observability features have been implemented:

### 1. Structured Logging

- Uses `zerolog` for JSON-structured logging
- Log levels: debug, info, warn, error, fatal
- Includes contextual information such as request IDs, durations, and component names
- Can be configured to output pretty-printed logs for development

### 2. Metrics Collection

- Uses Prometheus for metrics collection and storage
- Key metrics captured:
  - Request counts and latencies
  - Token usage (input and output)
  - Model performance (total latency, time to first token)
  - Error rates by type
  - Active request count
  - Memory usage
  - llama.cpp specific performance metrics

### 3. Tracing

- OpenTelemetry integration for distributed tracing
- Traces request flow from frontend to backend to model
- Captures spans for key operations

### 4. Visualization

- Grafana dashboard for metrics visualization
- Frontend metrics panel for quick insights
- Jaeger UI for trace exploration

### 5. Health Checks

- `/health` endpoint for basic health status
- `/readiness` endpoint for readiness checks
- Memory stats and uptime information

## Architecture

The observability stack consists of:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │ >>> │   Backend   │ >>> │ Model Runner│
│  (React/TS) │     │    (Go)     │     │ (Llama 3.2) │
└─────────────┘     └─────────────┘     └─────────────┘
                          v  v
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Grafana   │ <<< │ Prometheus  │     │   Jaeger    │
│ Dashboards  │     │  Metrics    │     │   Tracing   │
└─────────────┘     └─────────────┘     └─────────────┘
```

## Getting Started

### Environment Variables

The following environment variables can be set in `backend.env`:

```
# Observability configuration
LOG_LEVEL: info      # debug, info, warn, error
LOG_PRETTY: true     # Whether to output pretty-printed logs
TRACING_ENABLED: true  # Enable OpenTelemetry tracing
OTLP_ENDPOINT: jaeger:4318  # OpenTelemetry collector endpoint
```

### Accessing Dashboards

- **Metrics Dashboard**: http://localhost:3001 (Grafana)
- **Tracing UI**: http://localhost:16686 (Jaeger UI)
- **Prometheus**: http://localhost:9091

### Default Credentials

- **Grafana**: admin/admin

## Grafana Dashboards

AIWatch includes multiple pre-configured Grafana dashboards:

### 1. LLM Performance Dashboard (llm-dashboard.json)

This dashboard provides an overview of general LLM performance metrics:

- Model latency (p50 and p95)
- Time to first token
- Token usage by direction (input/output)
- API request counts
- Error counts
- Active request gauge

### 2. llama.cpp Performance Dashboard (llamacpp-dashboard.json)

This dashboard focuses specifically on llama.cpp performance metrics:

- Tokens per second gauge and time-series
- Memory per token usage 
- Context window size
- Prompt evaluation time (p50 and p95)
- Thread utilization
- Batch size
- Token processing over time

## Prometheus Metrics

AIWatch exposes the following metrics for Prometheus scraping:

### General Metrics

| Metric Name | Type | Description |
|-------------|------|-------------|
| `genai_app_http_requests_total` | Counter | HTTP request count by method, endpoint, status |
| `genai_app_http_request_duration_seconds` | Histogram | HTTP request duration |
| `genai_app_chat_tokens_total` | Counter | Token count by direction (input/output) |
| `genai_app_model_latency_seconds` | Histogram | Model response time |
| `genai_app_first_token_latency_seconds` | Histogram | Time to first token |
| `genai_app_errors_total` | Counter | Error count by type |
| `genai_app_active_requests` | Gauge | Number of currently active requests |
| `genai_app_model_memory_bytes` | Gauge | Model memory usage |

### llama.cpp Specific Metrics

| Metric Name | Type | Description |
|-------------|------|-------------|
| `genai_app_llamacpp_context_size` | Gauge | Context window size in tokens |
| `genai_app_llamacpp_prompt_eval_seconds` | Histogram | Time spent evaluating prompt |
| `genai_app_llamacpp_tokens_per_second` | Gauge | Token generation speed |
| `genai_app_llamacpp_memory_per_token_bytes` | Gauge | Memory usage per token |
| `genai_app_llamacpp_threads_used` | Gauge | Number of threads used for inference |
| `genai_app_llamacpp_batch_size` | Gauge | Batch size used for inference |

## Custom Metric Endpoints

In addition to standard Prometheus metrics, the application exposes:

- `/metrics/summary` - High-level metrics summary for the frontend
- `/metrics/log` - Endpoint to log metrics from the frontend
- `/metrics/error` - Endpoint to log errors from the frontend

## Logging

AIWatch uses structured JSON logging with the following features:

1. **Contextual Information**: Each log entry includes important context:
   - Request ID for request correlation
   - Component name for log source identification
   - User information when available
   - Duration for performance tracking

2. **Log Levels**: Five logging levels are available:
   - `debug`: Detailed information for debugging
   - `info`: General operational information
   - `warn`: Warnings that don't affect operation
   - `error`: Errors that affect specific operations
   - `fatal`: Critical errors requiring immediate attention

3. **Pretty Printing**: For development, logs can be formatted in a human-readable format.

4. **Request Logging**: HTTP requests are automatically logged with:
   - Request method and path
   - Response status code
   - Duration
   - User agent

## Tracing

AIWatch implements distributed tracing using OpenTelemetry:

1. **Span Context**: Traces follow requests through the system:
   - Frontend to Backend
   - Backend to Model
   - Between internal components

2. **Span Attributes**: Each span includes important attributes:
   - Request parameters
   - Component information
   - Correlation IDs

3. **Jaeger UI**: Traces can be visualized and analyzed in Jaeger.

## llama.cpp Metrics Details

The AIWatch platform collects and exposes metrics that specifically help monitor llama.cpp model performance:

### Tokens per Second

The number of tokens generated per second, a key performance metric for generation speed:

- Higher values indicate better performance
- Affected by model size, quantization, and hardware
- Tracked in real-time during token generation

### Context Window Size

The maximum number of tokens the model can process in a single context:

- Larger windows allow for longer conversations and more context
- Impacts memory usage significantly
- Fixed per model at initialization time

### Prompt Evaluation Time

The time spent processing the input prompt before generating the first token:

- Critical metric for responsiveness
- Increases with prompt length
- Can be a bottleneck for large prompts

### Memory per Token

The amount of memory used per token, a measure of memory efficiency:

- Lower values indicate better memory optimization
- Varies by model size and quantization level
- Important for resource planning and allocation

### Thread Utilization

The number of CPU threads used during inference:

- More threads can improve performance to a point
- Optimal value depends on hardware and model size
- Too many threads can cause performance degradation

### Batch Size

The number of tokens processed in a single batch:

- Larger batch sizes can improve throughput
- Impacts latency and memory usage
- Optimal value depends on hardware and model characteristics

## Performance Optimization

You can use the collected metrics to optimize llama.cpp model performance:

1. **Thread Count Tuning**: Adjust the number of threads based on the `genai_app_llamacpp_threads_used` metric and CPU utilization.

2. **Batch Size Optimization**: Modify batch size based on the `genai_app_llamacpp_tokens_per_second` and latency metrics.

3. **Memory Management**: Monitor `genai_app_llamacpp_memory_per_token_bytes` to ensure efficient memory usage.

4. **Prompt Processing**: Track `genai_app_llamacpp_prompt_eval_seconds` to identify potential bottlenecks in prompt handling.

5. **Context Size Adjustment**: Configure the context window size based on your application needs and available memory.

## FAQs

### How do I add a new metric?

Add new metrics to the `pkg/metrics/metrics.go` file following the existing patterns.

### How do I add distributed tracing to a new endpoint?

Use the `tracing.StartSpan()` function to create a new span, and call the cleanup function when done:

```go
ctx, endSpan := tracing.StartSpan(r.Context(), "operation_name")
defer endSpan()

// Your code here
```

### How do I customize the Grafana dashboard?

The dashboard is defined in `grafana/provisioning/dashboards/llm-dashboard.json` and `grafana/provisioning/dashboards/llamacpp-dashboard.json`. You can edit them directly or export a new version from the Grafana UI.

### How can I extend the llama.cpp metrics?

To add new llama.cpp metrics:

1. Define new metrics in `pkg/metrics/metrics.go`
2. Collect the metrics in your code
3. Call `RecordLlamaCppMetrics()` with the new values
4. Update the Grafana dashboard to display the new metrics

### How do I troubleshoot missing metrics?

1. Check that Prometheus is running and can reach the backend
2. Verify that the metrics endpoint is exposed at `/metrics`
3. Ensure that the model you're using is indeed a llama.cpp model
4. Check log files for any errors related to metrics collection
