package auth

import (
	"time"

	"github.com/BelyaevEI/microservices_auth/internal/cache"
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	"github.com/BelyaevEI/platform_common/pkg/db"
)

type serv struct {
	authRepository         repository.AuthRepository
	cache                  cache.UserCache
	txManager              db.TxManager
	refreshTokenSecretKey  string
	refreshTokenExpiration time.Duration
	accessTokenSecretKey   string
	accessTokenExpiration  time.Duration
}

// NewService creates a new auth service.
func NewService(
	authRepository repository.AuthRepository,
	cache cache.UserCache,
	txManager db.TxManager,
	refreshTokenSecretKey string,
	refreshTokenExpiration time.Duration,
	accessTokenSecretKey string,
	accessTokenExpiration time.Duration,
) service.AuthService {
	return &serv{
		authRepository:         authRepository,
		cache:                  cache,
		txManager:              txManager,
		refreshTokenSecretKey:  refreshTokenSecretKey,
		refreshTokenExpiration: refreshTokenExpiration,
		accessTokenSecretKey:   accessTokenSecretKey,
		accessTokenExpiration:  accessTokenExpiration,
	}
}
