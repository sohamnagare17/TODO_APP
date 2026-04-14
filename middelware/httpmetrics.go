package middleware

import (
	"go-sqlite/metrics"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// wrap response writer to capture status
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// call next handler
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		method := r.Method
		route := r.URL.Path
		status := http.StatusText(rw.statusCode)

		// metrics update
		metrics.HttpRequestsTotal.WithLabelValues(method, route, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, route).Observe(duration)

		if rw.statusCode >= 400 {
			metrics.HttpErrorsTotal.WithLabelValues(method, route).Inc()
		}
	})
}
