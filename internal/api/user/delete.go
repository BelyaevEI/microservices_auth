package user

import (
	"context"

	desc "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUserByID deletes a user.
func (i *Implementation) DeleteUserByID(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.DeleteUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
