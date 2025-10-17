package telemetry

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/metric"
)

var (
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rpcv2_hist_request_duration_seconds",
		Help: "Duration of RPC requests in seconds",
	}, []string{"method", "status"})

	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rpcv2_hist_requests_total",
		Help: "Total number of RPC requests",
	}, []string{"method", "status"})
)

type Metrics struct {
	requestDuration prometheus.ObserverVec
	requestCount    *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		requestDuration: requestDuration,
		requestCount:    requestCount,
	}
}

func (m *Metrics) RecordRequest(method string, status string, duration float64) {
	m.requestDuration.WithLabelValues(method, status).Observe(duration)
	m.requestCount.WithLabelValues(method, status).Inc()
}