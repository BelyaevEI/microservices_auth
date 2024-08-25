package user

import (
	"context"
	"log"

	"github.com/BelyaevEI/microservices_auth/internal/converter"
	desc "github.com/BelyaevEI/microservices_auth/pkg/user_v1"
)

// GetUserByID gets a user.
func (i *Implementation) GetUserByID(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %d, created_at: %v, updated_at: %v\n", user.ID, user.Info.Name, user.Info.Email, user.Info.Role, user.CreatedAt, user.UpdatedAt)

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
