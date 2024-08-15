package user

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
)

func (s serv) GetUserByID(ctx context.Context, id int64) (*model.User, error) {

	user, err := s.cache.GetUserByID(ctx, id)
	if err != nil {
		user, err = s.userRepository.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
