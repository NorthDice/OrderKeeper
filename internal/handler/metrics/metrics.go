package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)

	databaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	databaseConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	cacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	cacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	ordersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_total",
			Help: "Total number of orders",
		},
		[]string{"status"},
	)

	userRegistrationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}
		start := time.Now()

		activeConnections.Inc()
		defer activeConnections.Dec()

		c.Next()

		duration := time.Since(start)
		statusCode := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			statusCode,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			statusCode,
		).Observe(duration.Seconds())
	}
}

func RecordCacheHit(cacheType string) {
	cacheHitsTotal.WithLabelValues(cacheType).Inc()
}

func RecordCacheMiss(cacheType string) {
	cacheMissesTotal.WithLabelValues(cacheType).Inc()
}

func RecordOrder(status string) {
	ordersTotal.WithLabelValues(status).Inc()
}

func RecordUserRegistration() {
	userRegistrationsTotal.Inc()
}

func UpdateDatabaseConnections(active, idle int) {
	databaseConnectionsActive.Set(float64(active))
	databaseConnectionsIdle.Set(float64(idle))
}
