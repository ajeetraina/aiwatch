package middleware

import (
	"net/http"

	"github.com/ajeetraina/aiwatch/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
)

// TracingMiddleware adds OpenTelemetry tracing to HTTP requests
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start a new span for this request
		ctx, span := tracing.StartSpan(r.Context(), "http_request")
		defer span.End()

		// Add some attributes to the span
		tracing.AddAttributes(ctx,
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.user_agent", r.UserAgent()),
		)

		// Create a custom response writer to capture the status code
		writer := &responseWriter{w, http.StatusOK}

		// Use the context with the span
		next.ServeHTTP(writer, r.WithContext(ctx))

		// Add the response status code to the span
		tracing.AddAttribute(ctx, "http.status_code", writer.status)
	})
}
