package env

import (
	"errors"
	"net"
	"os"
)

const (
	prometheusHostEnv = "PROMETHEUS_HOST"
	prometheusPortEnv = "PROMETHEUS_PORT"
	prometheusPathEnv = "PROMETHEUS_PATH"
)

type prometheusConfig struct {
	prometheusHost string
	prometheusPort string
	prometheusPath string
}

// NewPrometheusConfig - конфигурация для сервера prometheus
func NewPrometheusConfig() (*prometheusConfig, error) {
	prometheusHost := os.Getenv(prometheusHostEnv)
	if len(prometheusHost) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	prometheusPort := os.Getenv(prometheusPortEnv)
	if len(prometheusPort) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	prometheusPath := os.Getenv(prometheusPathEnv)
	if len(prometheusPath) == 0 {
		return nil, errors.New("prometheus path not found")
	}

	return &prometheusConfig{
		prometheusHost: prometheusHost,
		prometheusPort: prometheusPort,
		prometheusPath: prometheusPath,
	}, nil
}

// Address - адрес сервера
func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.prometheusHost, cfg.prometheusPort)
}

// Path
func (cfg *prometheusConfig) Path() string {
	return "/" + cfg.prometheusPath
}
