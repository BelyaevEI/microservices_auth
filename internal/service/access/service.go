package access

import (
	"github.com/BelyaevEI/microservices_auth/internal/cache"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	"github.com/BelyaevEI/platform_common/pkg/db"
)

type serv struct {
	cache                cache.UserCache
	txManager            db.TxManager
	authPrefix           string
	accessTokenSecretKey string
}

// NewService creates a new access service.
func NewService(
	cache cache.UserCache,
	txManager db.TxManager,
	authPrefix string,
	accessTokenSecretKey string,
) service.AccessService {
	return &serv{
		cache:                cache,
		txManager:            txManager,
		authPrefix:           authPrefix,
		accessTokenSecretKey: accessTokenSecretKey,
	}
}
