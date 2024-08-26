package auth

import (
	"context"
	"github.com/BelyaevEI/microservices_auth/internal/service"
	"github.com/BelyaevEI/microservices_auth/pkg/access_v1"
	desc "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Implementation represents a auth API implementation.
type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// Check implements access_v1.AccessV1Server.
func (s *Implementation) Check(context.Context, *access_v1.CheckRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedAccessV1Server implements access_v1.AccessV1Server.
func (s *Implementation) mustEmbedUnimplementedAccessV1Server() {
	panic("unimplemented")
}

// NewImplementation creates a new auth API implementation.
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
