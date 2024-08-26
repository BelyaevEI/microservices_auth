package app

import (
	"context"
	"log"

	"github.com/BelyaevEI/microservices_auth/internal/api/access"
	"github.com/BelyaevEI/microservices_auth/internal/api/auth"
	"github.com/BelyaevEI/microservices_auth/internal/api/user"
	"github.com/BelyaevEI/microservices_auth/internal/cache"
	"github.com/BelyaevEI/microservices_auth/internal/client/kafka"
	"github.com/BelyaevEI/platform_common/pkg/cache/redis"
	"github.com/BelyaevEI/platform_common/pkg/closer"
	"github.com/BelyaevEI/platform_common/pkg/db"
	"github.com/BelyaevEI/platform_common/pkg/db/pg"
	"github.com/BelyaevEI/platform_common/pkg/db/transaction"

	userCache "github.com/BelyaevEI/microservices_auth/internal/cache/user"
	kafkaConsumer "github.com/BelyaevEI/microservices_auth/internal/client/kafka/consumer"
	"github.com/BelyaevEI/microservices_auth/internal/config"
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	authRepository "github.com/BelyaevEI/microservices_auth/internal/repository/auth"
	userRepository "github.com/BelyaevEI/microservices_auth/internal/repository/user"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	accessService "github.com/BelyaevEI/microservices_auth/internal/service/access"
	authService "github.com/BelyaevEI/microservices_auth/internal/service/auth"
	"github.com/BelyaevEI/microservices_auth/internal/service/consumer"
	userSaverConsumer "github.com/BelyaevEI/microservices_auth/internal/service/consumer/user_saver"
	userService "github.com/BelyaevEI/microservices_auth/internal/service/user"
	cacheClient "github.com/BelyaevEI/platform_common/pkg/cache"

	"github.com/IBM/sarama"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	redisConfig         config.RedisConfig
	kafkaConsumerConfig config.KafkaConsumerConfig
	jwtConfig           config.JWTConfig

	dbClient       db.Client
	redisPool      *redigo.Pool
	redisClient    cacheClient.Client
	cache          cache.UserCache
	txManager      db.TxManager
	userRepository repository.UserRepository
	authRepository repository.AuthRepository
	userService    service.UserService
	authService    service.AuthService
	accessService  service.AccessService
	userImpl       *user.Implementation
	authImpl       *auth.Implementation
	accessImpl     *access.Implementation

	userSaverConsumer    consumer.Servicer
	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %s", err.Error())
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

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

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
			s.Cache(),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthRepository(ctx),
			s.Cache(),
			s.TxManager(ctx),
			s.JWTConfig().RefreshSecretKey(),
			s.JWTConfig().RefreshExpiration(),
			s.JWTConfig().AccessSecretKey(),
			s.JWTConfig().AccessExpiration(),
		)
	}

	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(
			s.Cache(),
			s.TxManager(ctx),
			s.JWTConfig().AuthPrefix(),
			s.JWTConfig().AccessSecretKey(),
		)
	}

	return s.accessService
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

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

func (s *serviceProvider) RedisClient() cacheClient.Client {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig().ConnectionTimeout())
	}

	return s.redisClient
}

func (s *serviceProvider) Cache() cache.UserCache {
	if s.cache == nil {
		s.cache = userCache.NewCache(s.RedisClient())
	}
	return s.cache
}

func (s *serviceProvider) UserSaverConsumer(ctx context.Context) consumer.Servicer {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = userSaverConsumer.NewService(
			s.UserRepository(ctx),
			s.Consumer(),
		)
	}

	return s.userSaverConsumer
}

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

func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := config.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}
