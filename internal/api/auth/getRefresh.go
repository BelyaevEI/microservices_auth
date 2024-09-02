package auth

import (
	"context"
	"errors"

	errorAuth "github.com/BelyaevEI/microservices_auth/internal/errors"
	desc "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken implements auth.GetRefreshToken
func (s *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := s.authService.GetRefreshToken(ctx, req.RefreshToken)
	if errors.Is(err, errorAuth.ErrInvalidRefreshToken) {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
