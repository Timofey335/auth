package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/Timofey335/platform_common/pkg/closer"
	"github.com/Timofey335/platform_common/pkg/db"
	"github.com/Timofey335/platform_common/pkg/db/pg"
	"github.com/Timofey335/platform_common/pkg/db/transaction"
	"github.com/Timofey335/platform_common/pkg/kafka"
	kafkaConsumer "github.com/Timofey335/platform_common/pkg/kafka/consumer"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/Timofey335/auth/internal/api/access"
	"github.com/Timofey335/auth/internal/api/user"
	cacheImplementation "github.com/Timofey335/auth/internal/cache"
	userCache "github.com/Timofey335/auth/internal/cache/user"
	"github.com/Timofey335/auth/internal/client/cache"
	"github.com/Timofey335/auth/internal/client/cache/redis"
	"github.com/Timofey335/auth/internal/config"
	"github.com/Timofey335/auth/internal/config/env"
	"github.com/Timofey335/auth/internal/repository"
	accessRepository "github.com/Timofey335/auth/internal/repository/access"
	userRepository "github.com/Timofey335/auth/internal/repository/user"
	"github.com/Timofey335/auth/internal/service"
	accessService "github.com/Timofey335/auth/internal/service/access"
	userSaverConsumer "github.com/Timofey335/auth/internal/service/consumer/user_saver"
	userService "github.com/Timofey335/auth/internal/service/user"
)

type serviceProvider struct {
	authConfig          config.AuthConfig
	grpcConfig          config.GRPCConfig
	httpConfig          config.HTTPConfig
	pgConfig            config.PGConfig
	kafkaConsumerConfig config.KafkaConsumerConfig
	redisConfig         config.RedisConfig
	swaggerConfig       config.SwaggerConfig

	dbClient  db.Client
	txManager db.TxManager

	accessRepository repository.AccessRepository
	userRepository   repository.UserRepository
	cache            cacheImplementation.UserCache

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	accessService     service.AccessService
	userService       service.UserService
	userSaverConsumer service.ConsumerService

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler

	servImplementation       *user.Implementation
	accessServImplementation *access.AccessImplementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// AuthConfig - получает секреты из env файла
func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := env.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config: %s", err.Error())
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

// PGConfig - инициализирует конфигурацию БД из env файла
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig - инициализирует конфигурацию GRPC сервера
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// HTTPConfig - инициализирует конфигурацию HTTP сервера
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}
		s.httpConfig = cfg
	}

	return s.httpConfig
}

// SwaggerConfig - инициализирует конфигурацию Swagger сервера
func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}
		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// RedisConfig - инициализирует конфигурацию redis
func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

// RedisConfig - инициализирует конфигурацию kafka
func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

// RedisPool - конфигурация redis pool
func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

// RedisClient - создает клиент redis
func (s *serviceProvider) RedisClient(ctx context.Context) cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	err := s.redisClient.Ping(ctx)
	if err != nil {
		log.Fatalf("failed connect to redis: %v", err)
	}

	return s.redisClient
}

// Cache - имплементация кэш сервиса
func (s *serviceProvider) Cache(ctx context.Context) cacheImplementation.UserCache {
	if s.cache == nil {
		s.cache = userCache.NewCache(s.RedisClient(ctx))
	}

	return s.cache
}

// DBClient - инициализирует подключение к БД
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager - инициализирует Transaction Manager
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// AccessService - инициализация сервисного слоя access
func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(
			s.AccessRepository(ctx),
			s.UserRepository(ctx),
			s.Cache(ctx),
			s.TxManager(ctx),
			s.AuthConfig(),
		)
	}

	return s.accessService
}

// AccessRepository - инициализация repo слоя Access
func (s *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepository.NewRepository(s.DBClient(ctx))
	}

	return s.accessRepository
}

// UserRepository - инициализация repo слоя User
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService - инициализация сервисного слоя user
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.Cache(ctx),
			s.TxManager(ctx),
			s.AuthConfig(),
		)
	}

	return s.userService
}

// UserSaverConsumer - инийиализирует userSaverConsumer
func (s *serviceProvider) UserSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = userSaverConsumer.NewService(
			s.UserRepository(ctx),
			s.Consumer(),
		)
	}

	return s.userSaverConsumer
}

// Consumer - инициализирует consumer
func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

// ConsumerGroup - инициализирует consumerGroup
func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

// ConsumerGroupHandler - инициализирует consumerGroupHandler
func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}

// ServImplementation - инициализация api слоя
func (s *serviceProvider) ServImplementation(ctx context.Context) *user.Implementation {
	if s.servImplementation == nil {
		s.servImplementation = user.NewImplementation(s.UserService(ctx))
	}

	return s.servImplementation
}

// AccessServImplementation - инициализация api слоя access
func (s *serviceProvider) AccessServImplementation(ctx context.Context) *access.AccessImplementation {
	if s.accessServImplementation == nil {
		s.accessServImplementation = access.NewAccessImplementation(s.AccessService(ctx))
	}

	return s.accessServImplementation
}
