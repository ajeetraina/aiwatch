package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsMiddleware adds metrics collection middleware
func MetricsMiddleware(requestCounter *prometheus.CounterVec, requestDuration *prometheus.HistogramVec, activeRequests prometheus.Gauge) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create a custom response writer to capture the status code
			writer := &responseWriter{w, http.StatusOK}
			
			// Increment active requests counter
			activeRequests.Inc()
			
			// Call the next handler
			next.ServeHTTP(writer, r)
			
			// Decrement active requests counter
			activeRequests.Dec()
			
			// Calculate request duration
			duration := time.Since(start)
			
			// Record metrics
			requestCounter.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(writer.status)).Inc()
			requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
		})
	}
}
