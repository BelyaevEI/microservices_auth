package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

func (s serv) CreateUser(ctx context.Context, userCreate *model.UserCreate) (int64, error) {

	id, err := s.userRepository.CreateUser(ctx, userCreate)
	if err != nil {
		return 0, err
	}

	err = s.cache.CreateUser(ctx, &model.User{ID: id, Info: model.UserInfo{Name: userCreate.Name,
		Email: userCreate.Email,
		Role:  userCreate.Role}})
	if err != nil {
		return 0, nil
	}

	return id, nil
}
