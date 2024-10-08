package auth

import (
	"context"

	"github.com/pkg/errors"

	errorAuth "github.com/BelyaevEI/microservices_auth/internal/errors"
	desc "github.com/BelyaevEI/microservices_auth/pkg/auth_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccessToken implements auth.GetAccessToken
func (s *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := s.authService.GetAccessToken(ctx, req.RefreshToken)
	if errors.Is(err, errorAuth.ErrInvalidAccessToken) {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
