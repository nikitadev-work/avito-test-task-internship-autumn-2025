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
	businessOpsTotal    *prometheus.CounterVec
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

		businessOpsTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "business_operations_total",
				Help: "Count of business-level operations.",
			},
			[]string{"service", "operation", "result"},
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

		labels := prometheus.Labels{
			"service": serviceName,
			"method":  r.Method,
			"path":    r.URL.Path,
			"status":  strconv.Itoa(rec.status),
		}

		httpRequestsTotal.With(labels).Inc()
		httpRequestDuration.With(labels).Observe(duration.Seconds())
	})
}

func IncBusinessOperation(service, operation, result string) {
	if businessOpsTotal == nil {
		return
	}
	businessOpsTotal.With(prometheus.Labels{
		"service":   service,
		"operation": operation,
		"result":    result,
	}).Inc()
}
