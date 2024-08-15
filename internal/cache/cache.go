package cache

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

type UserCache interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	DeleteUserByID(ctx context.Context, id int64) error
}
