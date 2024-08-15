package user

import (
	"context"
	"errors"
	"strconv"

	"github.com/BelyaevEI/microservices_auth/internal/cache/user/converter"
	modelCache "github.com/BelyaevEI/microservices_auth/internal/cache/user/model"
	"github.com/BelyaevEI/microservices_auth/internal/model"
	redigo "github.com/gomodule/redigo/redis"
)

func (c *cacheImplementation) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := c.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, errors.New("user not found")
	}

	var user modelCache.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromCache(&user), nil
}
