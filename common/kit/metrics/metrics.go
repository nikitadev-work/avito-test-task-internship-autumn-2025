package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	metricsInitOnce     sync.Once
)

func InitMetrics() {
	metricsInitOnce.Do(func() {
		httpRequestsTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"service", "method", "path", "status"},
		)

		httpRequestDuration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method", "path", "status"},
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Middleware to record metrics
func HTTPMiddleware(serviceName string, next http.Handler) http.Handler {
	InitMetrics()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		if httpRequestsTotal == nil || httpRequestDuration == nil {
			return
		}

		path := r.URL.Path
		method := r.Method
		statusStr := strconv.Itoa(rec.status)

		labels := prometheus.Labels{
			"service": serviceName,
			"method":  method,
			"path":    path,
			"status":  statusStr,
		}

		httpRequestsTotal.With(labels).Inc()
		httpRequestDuration.With(labels).Observe(duration.Seconds())
	})
}
