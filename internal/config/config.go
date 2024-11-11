package config

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

// AuthConfig - интерфейс для auth
type AuthConfig interface {
	AccessTokenSecretKey() string
	AccessTokenExpiration() int64
	RefreshTokenSecretKey() string
	RefreshTokenExpiration() int64
}

// GRPCConfig - интерфейс конфигурации grpc сервера
type GRPCConfig interface {
	Address() string
}

// HTTPConfig - интерфейс конфигурации http сервера
type HTTPConfig interface {
	Address() string
}

// KafkaConsumerConfig - интерфейс kafka
type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

// Load - считывает переменные из env файла
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// LoggerConfig - интерфейс для логгера
type LoggerConfig interface {
	LogLevel() string
	LogFilename() string
	LogFileMaxSize() int
	LogFileMaxBackups() int
	LogFileMaxAge() int
}

// PGConfig - интерфейс с методом DSN
type PGConfig interface {
	DSN() string
}

// PrometheusConfig - интерфейс для prometheus
type PrometheusConfig interface {
	Address() string
	Path() string
}

// RedisConfig - интерфейс конфигурации redis
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	UserExpiration() int64
}

// SwaggerConfig - интерфейс конфигурации swagger сервера
type SwaggerConfig interface {
	Address() string
}
