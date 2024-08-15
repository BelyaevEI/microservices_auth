package user

import (
	"github.com/BelyaevEI/microservices_auth/internal/cache"
	cacheClient "github.com/BelyaevEI/platform_common/pkg/cache"
)

type cacheImplementation struct {
	cl cacheClient.Client
}

func NewCache(cl cacheClient.Client) cache.UserCache {
	return &cacheImplementation{
		cl: cl,
	}
}
