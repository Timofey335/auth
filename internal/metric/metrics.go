package metric

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "my_space"
	appName   = "auth"
)

// Metrics - используемые типы метрик prometheus
type Metrics struct {
	requestCounter  prometheus.Counter
	responseCounter *prometheus.CounterVec
}

var metrics *Metrics

// Init - инициализация метрик
func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_requests_total",
				Help:      "Количество запросов к серверу",
			},
		),

		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_responses_total",
				Help:      "Количество ответов от сервера",
			},
			[]string{"status", "method"},
		),
	}

	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}
