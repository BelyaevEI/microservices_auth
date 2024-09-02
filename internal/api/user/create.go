package user

import (
	"context"
	"log"

	"github.com/BelyaevEI/microservices_auth/internal/converter"
	desc "github.com/BelyaevEI/microservices_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser creates a new user.
func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Error(codes.InvalidArgument, "password and password confirm do not match")
	}

	id, err := i.userService.CreateUser(ctx, converter.ToUserCreateFromDesc(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
