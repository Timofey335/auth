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
	requestCounter prometheus.Counter
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
	}

	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
