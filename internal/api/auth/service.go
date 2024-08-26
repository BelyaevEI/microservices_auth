package auth

import (
	"github.com/BelyaevEI/microservices_auth/internal/service"
	desc "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"
)

// Implementation represents a auth API implementation.
type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation creates a new auth API implementation.
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
