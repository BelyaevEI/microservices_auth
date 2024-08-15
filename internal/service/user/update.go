package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

func (s serv) UpdateUserByID(ctx context.Context, id int64, userUpdate *model.UserUpdate) error {

	err := s.userRepository.UpdateUserByID(ctx, id, userUpdate)
	if err != nil {
		return err
	}

	err = s.cache.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.cache.CreateUser(ctx, &model.User{ID: id, Info: model.UserInfo{
		Name:  *userUpdate.Name,
		Email: *userUpdate.Email,
		Role:  userUpdate.Role}})
	if err != nil {
		return err
	}

	return nil
}
