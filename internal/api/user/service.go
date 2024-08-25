package user

import (
	"github.com/BelyaevEI/microservices_auth/internal/service"
	desc "github.com/BelyaevEI/microservices_auth/pkg/user_v1"
)

// Implementation represents a user API implementation.
type Implementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation creates a new user API implementation.
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
