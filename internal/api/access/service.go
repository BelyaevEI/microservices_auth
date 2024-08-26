package access

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/BelyaevEI/microservices_auth/internal/service"
	desc "github.com/BelyaevEI/microservices_auth/pkg/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implementation represents a access API implementation.
type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewImplementation creates a new access API implementation.
func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}

// Check implements access.Check
func (s *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*empty.Empty, error) {
	err := s.accessService.Check(ctx, req.EndpointAddress)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return nil, nil
}
