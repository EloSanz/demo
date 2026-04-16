package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

// MetricsHandler returns the Prometheus metrics endpoint.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// MetricsMiddleware tracks status code, duration and request count.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Usamos el responseWriter que ya tenías en middleware.go para capturar el status
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.status)
		
		// Registramos las métricas
		// Nota: En producción, idealmente "path" debería estar sanitizado para evitar explosión de cardinalidad
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
