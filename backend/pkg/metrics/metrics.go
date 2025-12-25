package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// WebSocket metrics
	wsConnectionsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_total",
			Help: "Total number of active WebSocket connections",
		},
	)

	wsPriceUpdatesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_price_updates_total",
			Help: "Total number of price updates sent via WebSocket",
		},
	)

	// Trading metrics
	ordersCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total number of orders created",
		},
		[]string{"side", "type"},
	)

	tradesExecuted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "trades_executed_total",
			Help: "Total number of trades executed",
		},
	)

	tradeVolume = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "trade_volume_total",
			Help: "Total trading volume in USD",
		},
		[]string{"symbol"},
	)

	// Cache metrics
	cacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)

// Middleware creates a Prometheus middleware for Fiber
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())
		path := c.Route().Path
		method := c.Method()

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}

// Handler returns the Prometheus metrics handler for Fiber
func Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
		handler(c.Context())
		return nil
	}
}

// Recording functions
func RecordWSConnection(delta int) {
	wsConnectionsTotal.Add(float64(delta))
}

func RecordPriceUpdate() {
	wsPriceUpdatesTotal.Inc()
}

func RecordOrderCreated(side, orderType string) {
	ordersCreated.WithLabelValues(side, orderType).Inc()
}

func RecordTradeExecuted() {
	tradesExecuted.Inc()
}

func RecordTradeVolume(symbol string, volume float64) {
	tradeVolume.WithLabelValues(symbol).Add(volume)
}

func RecordCacheHit() {
	cacheHits.Inc()
}

func RecordCacheMiss() {
	cacheMisses.Inc()
}
