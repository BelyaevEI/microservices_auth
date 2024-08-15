package app

import (
	"context"
	"log"

	"github.com/BelyaevEI/microservices_auth/internal/api/user"
	"github.com/BelyaevEI/microservices_auth/internal/cache"
	"github.com/BelyaevEI/platform_common/pkg/cache/redis"
	"github.com/BelyaevEI/platform_common/pkg/closer"
	"github.com/BelyaevEI/platform_common/pkg/db"
	"github.com/BelyaevEI/platform_common/pkg/db/pg"
	"github.com/BelyaevEI/platform_common/pkg/db/transaction"

	userCache "github.com/BelyaevEI/microservices_auth/internal/cache/user"
	"github.com/BelyaevEI/microservices_auth/internal/config"
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	userRepository "github.com/BelyaevEI/microservices_auth/internal/repository/user"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	userService "github.com/BelyaevEI/microservices_auth/internal/service/user"
	cacheClient "github.com/BelyaevEI/platform_common/pkg/cache"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig

	dbClient       db.Client
	redisPool      *redigo.Pool
	redisClient    cacheClient.Client
	cache          cache.UserCache
	txManager      db.TxManager
	userRepository repository.UserRepository
	userService    service.UserService
	userImpl       *user.Implementation
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
