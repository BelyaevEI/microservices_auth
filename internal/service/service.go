package service

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

// UserService represents a service for user entities.
type UserService interface {
	CreateUser(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	UpdateUserByID(ctx context.Context, id int64, userUpdate *model.UserUpdate) error
	DeleteUserByID(ctx context.Context, id int64) error
}

// AuthService represents a service for auth entities.
type AuthService interface {
	Login(ctx context.Context, userLogin *model.UserLogin) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// AccessService represents a service for access entities.
type AccessService interface {
	Check(ctx context.Context, endpointAddress string) error
}
