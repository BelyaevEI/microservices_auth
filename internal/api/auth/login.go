package auth

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/model"
	desc "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login implements auth.Login
func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := s.authService.Login(ctx, &model.UserLogin{
		ID:       req.Id,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
