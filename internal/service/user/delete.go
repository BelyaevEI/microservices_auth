package user

import (
	"context"
)

func (s serv) DeleteUserByID(ctx context.Context, id int64) error {

	err := s.userRepository.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.cache.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
